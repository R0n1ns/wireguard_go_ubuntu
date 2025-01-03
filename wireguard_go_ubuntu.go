package wireguard_go_ubuntu

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/telebot.v3"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// Структура для конфигурации пира
type PeerConfig struct {
	PublicKey  string `json:"public_key"`
	AllowedIPs string `json:"allowed_ips"`
	Endpoint   string `json:"endpoint"`
}

// Структура для клиента
type Client struct {
	Id               int        `json:"id"`
	Status           bool       `json:"status"`
	AddressClient    string     `json:"address_client"`
	PubkeyPath       string     `json:"pubkey_path"`
	PrivkeyPath      string     `json:"privkey_path"`
	PrivateClientKey string     `json:"private_client_key"`
	PublicClientKey  string     `json:"public_client_key"`
	Peer             PeerConfig `json:"peer"`
	PeerStr          string     `json:"peer_str"`
	Config           string     `json:"config"`
	TgId             int        `json:"tg_id"`
}

// Управление сервером WireGuard
type WireGuardConfig struct {
	PrivateKey string         `json:"private_key"`
	PublicKey  string         `json:"public_key"`
	Endpoint   string         `json:"endpoint"`
	ListenPort string         `json:"listen_port"`
	InterName  string         `json:"inter_name"`
	BotToken   string         `json:"bot_token"`
	Clients    map[int]Client `json:"clients"` // Используем карту клиентов
}

// ------------------------ сохранение и загрузка данных ------------------------
// Метод сохранения WireGuardConfig в JSON файл
func (config *WireGuardConfig) SaveToFile(filename string) error {
	// Преобразуем конфигурацию в JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	// Записываем данные в файл
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Метод загрузки WireGuardConfig из JSON файла
func (config *WireGuardConfig) LoadFromFile(filename string) error {
	// Проверяем, существует ли файл
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		// Файл не существует, ничего не делаем
		return nil
	}
	if err != nil {
		return err
	}

	// Читаем данные из файла
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Раскодируем JSON данные в структуру config
	err = json.Unmarshal(data, config)
	if err != nil {
		return err
	}
	return nil
}

// ------------------------ методы для клиентов ------------------------
// Остановка клиента
func (wg WireGuardConfig) StopClient(id int) {
	client, exists := wg.Clients[id]
	if !exists {
		log.Printf("Клиент с id %d не найден", id)
		return
	}
	defer func() { wg.Clients[id] = client }()
	//fmt.Println(wg.Clients[id])

	filePath := "/etc/wireguard/wg0.conf"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Ошибка чтения файла конфигурации: %v", err)
		return
	}
	defer restWireguard()

	fileContent := string(content)
	client.Status = false
	updatedContent := strings.Replace(fileContent, client.PeerStr, "", 1)

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Printf("Ошибка записи файла конфигурации: %v", err)
		return
	}

	log.Printf("Клиент с id %d остановлен", id)
}

// Активация клиента
func (wg WireGuardConfig) ActClient(id int) {
	client, exists := wg.Clients[id]
	if !exists {
		log.Printf("Клиент с id %d не найден", id)
		return
	}
	//fmt.Println(wg.Clients[id])
	defer restWireguard()

	filePath := "/etc/wireguard/wg0.conf"
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Ошибка чтения файла конфигурации: %v", err)
		return
	}
	defer func() { wg.Clients[id] = client }()

	client.Status = true
	updatedContent := string(content) + "\n" + client.PeerStr

	err = os.WriteFile(filePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Printf("Ошибка записи файла конфигурации: %v", err)
		return
	}

	log.Printf("Клиент с id %d активирован", id)
}

// Удаление клиента
func (wg *WireGuardConfig) DeleteClient(id int) {
	wg.StopClient(id)

	//err := os.Remove(fmt.Sprintf("/etc/wireguard/wg_client_%d_private", id))
	//if err != nil {
	//	log.Printf("Не удалось удалить файл: %v", err)
	//}
	//
	//err = os.Remove(fmt.Sprintf("/etc/wireguard/wg_client_%d_public", id))
	//if err != nil {
	//	log.Printf("Не удалось удалить файл: %v", err)
	//}

	delete(wg.Clients, id)
}

// вывод всех клиентов
func (clients *WireGuardConfig) AllClients() string {
	text := ""
	for id, client := range clients.Clients {
		var stat string
		if client.Status {
			stat = "Активен"
		} else {
			stat = "Остановлен"
		}
		text += fmt.Sprintf("Клиент %d статус %s адресс %s \n", id, stat, client.AddressClient)
	}
	return text

}

// Добавление клиента WireGuard
func (wg *WireGuardConfig) AddWireguardClient(clientID int) (Client, int, error) {
	// Инициализация карты клиентов, если она nil
	if wg.Clients == nil {
		wg.Clients = make(map[int]Client)
	}
	defer restWireguard()
	// Проверяем, существует ли клиент
	client, exists := wg.Clients[clientID]
	if !exists {
		client = Client{Id: clientID}
		wg.Clients[clientID] = client
	}
	defer func() { wg.Clients[clientID] = client }()
	// Генерация ключей для клиента
	var privateKey, publicKey bytes.Buffer
	cmd := exec.Command("wg", "genkey")
	cmd.Stdout = &privateKey
	err := cmd.Run()

	client.PrivateClientKey = strings.TrimSpace(privateKey.String())
	cmd = exec.Command("wg", "pubkey")
	cmd.Stdin = &privateKey
	cmd.Stdout = &publicKey
	err = cmd.Run()

	client.PublicClientKey = strings.TrimSpace(publicKey.String())
	client.AddressClient = fmt.Sprintf("10.0.0.%d/24", clientID)
	client.Peer.Endpoint = wg.Endpoint
	client.Peer.PublicKey = wg.PublicKey
	peer := fmt.Sprintf("\n[Peer]\nPublicKey = %s\nAllowedIPs = %s\n", strings.TrimSpace(publicKey.String()), fmt.Sprintf("10.0.0.%d/24", clientID))
	client.PeerStr = peer
	filePath := "/etc/wireguard/wg0.conf"
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return Client{}, 0, err
	}
	defer f.Close()
	if _, err = f.WriteString(peer); err != nil {
		return Client{}, 0, err

	}
	client.Status = true
	// Генерация и сохранение конфигурации клиента
	clientConfig := fmt.Sprintf(`[Interface]
Address = %s
PrivateKey = %s
DNS = 8.8.8.8

[Peer]
Endpoint = %s
PublicKey = %s
AllowedIPs = 0.0.0.0/0
    `, client.AddressClient, client.PrivateClientKey, wg.Endpoint, wg.PublicKey)

	client.Config = clientConfig
	return client, clientID, nil
}

// ------------------------ методы для сервера ------------------------
// автоматический запуск сервера wiregguard
func (wg *WireGuardConfig) Autostart() {
	wg.RandomPort()
	wg.GetIPAndInterfaceName()
	wg.GenServerKeys()
	wg.GenerateWireGuardConfig()
	// wg_client.CollectTraffic()
	wg.WireguardStart()
}

// генерируем ключи сервера
func (wg *WireGuardConfig) GenServerKeys() {
	//генерируем ключи
	var privateKey bytes.Buffer
	cmd := exec.Command("wg", "genkey")
	cmd.Stdout = &privateKey
	err := cmd.Run()
	if err != nil {
		log.Fatalf("failed to generate private key: %v", err)
	}
	// Сохраняем приватный ключ в переменную
	privatekey := strings.ReplaceAll(privateKey.String(), "\n", "")
	// Используем приватный ключ для генерации публичного ключа
	var publicKey bytes.Buffer
	cmd = exec.Command("wg", "pubkey")
	cmd.Stdin = &privateKey
	cmd.Stdout = &publicKey

	err = cmd.Run()
	if err != nil {
		log.Fatalf("failed to generate public key: %v", err)
	}
	publickkey := strings.ReplaceAll(publicKey.String(), "\n", "")
	//запись
	os.WriteFile("/etc/wireguard/privatekey", []byte(privatekey), 0600)
	os.WriteFile("/etc/wireguard/publickey", []byte(publickkey), 0600)
	// Сохраняем публичный ключ в переменную
	time.Sleep(time.Second * 5)
	wg.PublicKey = publickkey
	wg.PrivateKey = privatekey
}

// генерация рандомного порта
func (wg *WireGuardConfig) RandomPort() {
	wg.ListenPort = strconv.Itoa(rand.Intn(10000))
	//fmt.Println(wg.ListenPort)
}

// получение сети
type NetworkInterface struct {
	Name string
	IsUp bool
	IPs  []string
}

func (cfg *WireGuardConfig) GetIPAndInterfaceName() error {
	interfaces, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range interfaces {
		// Пропускаем неактивные интерфейсы или интерфейсы без нужных флагов
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return err
		}

		for _, addr := range addrs {
			var ip net.IP

			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Пропускаем не IPv4 адреса
			if ip == nil || ip.To4() == nil {
				continue
			}

			cfg.InterName = iface.Name
			cfg.Endpoint = ip.String() + ":" + cfg.ListenPort
			return nil
		}
	}

	return fmt.Errorf("не удалось найти подходящий IP-адрес и интерфейс")
}

// Функция для определения, является ли интерфейс проводным
func isWiredInterface(name string) bool {
	return strings.HasPrefix(name, "e") || strings.Contains(name, "eth") || strings.Contains(name, "en")
}

// Функция для определения, является ли интерфейс беспроводным
func isWirelessInterface(name string) bool {
	return strings.HasPrefix(name, "w") || strings.Contains(name, "wl") || strings.Contains(name, "wlan")
}

// Генерация конфигурации WireGuard
func (wg *WireGuardConfig) GenerateWireGuardConfig() {
	//генерация конфига для мурвера
	tmpl := `[Interface]
PrivateKey = {{.PrivateKey}}
Address = 10.0.0.1/24
ListenPort = {{.ListenPort}}
PostUp = iptables -A FORWARD -i %i -j ACCEPT; iptables -t nat -A POSTROUTING -o {{.InterName}} -j MASQUERADE
PostDown = iptables -D FORWARD -i %i -j ACCEPT; iptables -t nat -D POSTROUTING -o {{.InterName}} -j MASQUERADE`

	t := template.Must(template.New("wgConfig").Parse(tmpl))

	var buf bytes.Buffer
	if err := t.Execute(&buf, wg); err != nil {
	}
	// Генерация случайного числа от 0 до 99999
	// Открытие файла для записи
	filePath := "/etc/wireguard/wg0.conf"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Ошибка открытия файла: %v\n", err)
		return
	}
	defer file.Close()

	// Запись данных в файл
	_, err = file.Write(buf.Bytes())
	if err != nil {
		fmt.Printf("Ошибка записи в файл: %v\n", err)
		return
	}

	//log.Println("Конфигурация wireguard сгенерирована")

}

func (wg *WireGuardConfig) WireguardStart() {
	port := wg.ListenPort
	// настройка форвардинг
	file, err := os.OpenFile("/etc/sysctl.conf", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open /etc/sysctl.conf: %v", err)
	}
	defer file.Close()

	// Записываем строку в файл
	_, err = file.WriteString("net.ipv4.ip_forward=1\n")
	if err != nil {
		log.Fatalf("failed to write to /etc/sysctl.conf: %v", err)
	}
	prt := fmt.Sprintf("%s/udp", port)
	// Выполняем команду `sysctl -p` для применения изменений
	cmd := exec.Command("ufw", "allow", prt)
	err = cmd.Run()
	if err != nil {
		log.Printf("failed to apply sysctl changes: %v", err.Error())
	}
	// Выполняем команду `sysctl -p` для применения изменений
	cmd = exec.Command("sysctl", "-p")
	err = cmd.Run()
	if err != nil {
		log.Printf("failed to apply sysctl changes: %v", err.Error())
	}
	//включсение wireguard
	cmd = exec.Command("systemctl", "enable", "wg-quick@wg0.service")
	cmd.Run()
	err = cmd.Err
	if err != nil {
		log.Printf("failed to create keys : %v", err.Error())
	}
	//старт wireguard
	cmd = exec.Command("systemctl", "start", "wg-quick@wg0.service")
	cmd.Run()
	err = cmd.Err
	if err != nil {
		log.Printf("failed to create keys : %v", err.Error())
	}
	//log.Printf("Соединение wireguard запущено")
}
func restWireguard() {
	cmd := exec.Command("systemctl", "restart", "wg-quick@wg0")
	err := cmd.Run()
	if err != nil {
		log.Printf("failed to apply sysctl changes: %v", err.Error())
	}

}

// Отправка конфигурации через Telegram
func (wg *WireGuardConfig) SendConfigToUserTg(user_id int) {
	if wg.BotToken == "" {
		var value string
		fmt.Print("Пожалуйста введите токен бота,или 0 для отмены")
		fmt.Scanln(&value)
		if value != "0" {
			wg.BotToken = value
		} else {
			return
		}
	}
	//создание бота
	Cl, _ := wg.Clients[user_id]

	chatID := telebot.ChatID(int64(Cl.TgId))
	bot, err := telebot.NewBot(telebot.Settings{
		Token: wg.BotToken,
	})

	//файл с конфигураций
	reader := strings.NewReader(Cl.Config)
	// Создаем документ для отправки, передавая reader как содержимое файла
	document := &telebot.Document{
		File:     telebot.FromReader(reader), // Используем io.Reader
		FileName: "wgconf.conf",              // Указываем имя файла
		Caption:  "WireGuard Configuration",  // Опциональная подпись к файлу
	}
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}
	//отправка файла
	if _, err := bot.Send(chatID, document); err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
}

// удаление wireguard
func (wg *WireGuardConfig) DropWireguard() {
	// очистка папки
	cmd := exec.Command("rm", "-rf", "/etc/wireguard/*")
	cmd.Run()
	err := cmd.Err
	if err != nil {
		log.Printf("failed to create keys : %v", err.Error())
	}
	//отключение wireguard
	cmd = exec.Command("systemctl", "disable", "wg-quick@wg0.service")
	cmd.Run()
	err = cmd.Err
	if err != nil {
		log.Printf("failed to create keys : %v", err.Error())
	}
	filePath := "/etc/sysctl.conf"

	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	// Читаем содержимое файла построчно
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Добавляем в список все строки, кроме той, которую нужно удалить
		if !strings.Contains(line, "net.ipv4.ip_forward=1") {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	// Перезаписываем файл без нужной строки
	err = os.WriteFile(filePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		log.Fatalf("failed to write file: %v", err)
	}

	fmt.Println("Line removed successfully")
	log.Printf("Папка конфиураций wireguuard очищена")
}

// // Сбор трафика
//
//	func (wg *WireGuardConfig) CollectTraffic() {
//		cmd := exec.Command("wg-json")
//		go cmd.Run() // Запускаем в горутине
//		log.Println("Сбор трафика начат. Для остановки используйте Ctrl+C.")
//	}
//
// Структура для хранения выходных данных команды wg-json

type PeerStats struct {
	TransferRx uint64
	TransferTx uint64
}

type PeerTraffic struct {
	TrafficRx uint64 `json:"traffic_rx"`
	TrafficTx uint64 `json:"traffic_tx"`
}

// Сбор трафика, возвращает map[id]PeerTraffic
func (wg *WireGuardConfig) CollectTraffic() (map[string]PeerTraffic, error) {
	cmd := exec.Command("wg", "show")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to execute wg show: %v", err)
	}

	output := out.String()
	lines := strings.Split(output, "\n")

	trafficData := make(map[string]PeerTraffic)

	var currentPeer string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "peer:") {
			currentPeer = strings.TrimSpace(strings.TrimPrefix(line, "peer:"))
		} else if strings.HasPrefix(line, "transfer:") && currentPeer != "" {
			// Пример: "transfer: 3.48 MiB received, 33.46 MiB sent"
			transferParts := strings.Split(line, ",")
			if len(transferParts) != 2 {
				continue
			}

			rxStr := strings.TrimSpace(strings.TrimPrefix(transferParts[0], "transfer:"))
			txStr := strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(transferParts[1], ""), "sent"))

			rxBytes, err := parseTraffic(rxStr)
			if err != nil {
				log.Printf("Error parsing transfer Rx for peer %s: %v", currentPeer, err)
				continue
			}

			txBytes, err := parseTraffic(txStr)
			if err != nil {
				log.Printf("Error parsing transfer Tx for peer %s: %v", currentPeer, err)
				continue
			}

			trafficData[currentPeer] = PeerTraffic{
				TrafficRx: rxBytes,
				TrafficTx: txBytes,
			}
		}
	}

	return trafficData, nil
}

// parseTraffic преобразует строку трафика (например, "3.48 MiB") в байты
func parseTraffic(trafficStr string) (uint64, error) {
	parts := strings.Fields(trafficStr)
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid traffic format: %s", trafficStr)
	}

	value, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse traffic value: %v", err)
	}

	unit := parts[1]
	switch unit {
	case "B":
		return uint64(value), nil
	case "KiB":
		return uint64(value * 1024), nil
	case "MiB":
		return uint64(value * 1024 * 1024), nil
	case "GiB":
		return uint64(value * 1024 * 1024 * 1024), nil
	default:
		return 0, fmt.Errorf("unknown traffic unit: %s", unit)
	}
}
