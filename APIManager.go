package wireguard_go_ubuntu

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

var cl = WireGuardConfig{}

// Обработчики API
// Добавление клиента
func AddClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := map[string]int{"id": 0}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &id); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	client, clID, err := cl.AddWireguardClient(id["id"])
	if err != nil {
		responseError(w, "Error adding client", http.StatusInternalServerError)
		return
	}

	responseJSON(w, struct {
		Client Client `json:"client"`
		ID     int    `json:"id"`
	}{Client: client, ID: clID})
}

// Удаление клиента
func DeleteClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := map[string]int{"id": 0}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &id); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	cl.DeleteClient(id["id"])
	responseJSON(w, map[string]string{"status": "Client deleted"})
}

// Список всех клиентов
func GetAllClientsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	clients := cl.AllClients()
	responseJSON(w, map[string]string{"clients": clients})
}

// Активация клиента
func ActivateClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := map[string]int{"id": 0}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &id); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	cl.ActClient(id["id"])
	responseJSON(w, map[string]string{"status": "Client activated"})
}

// Остановка клиента
func StopClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	id := map[string]int{"id": 0}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := json.Unmarshal(data, &id); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	cl.StopClient(id["id"])
	responseJSON(w, map[string]string{"status": "Client stopped"})
}

// Старт сервера WireGuard
func StartServerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	cl.Autostart()
	responseJSON(w, map[string]string{"status": "Server started"})
}

// Общие функции для ответа JSON
func responseJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func responseError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func main() {
	http.HandleFunc("/addClient", AddClientHandler)
	http.HandleFunc("/deleteClient", DeleteClientHandler)
	http.HandleFunc("/getAllClients", GetAllClientsHandler)
	http.HandleFunc("/activateClient", ActivateClientHandler)
	http.HandleFunc("/stopClient", StopClientHandler)
	http.HandleFunc("/startServer", StartServerHandler)

	log.Println("API server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
