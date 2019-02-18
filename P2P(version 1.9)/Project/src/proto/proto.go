package proto

import "encoding/json"

// Request -- запрос клиента к серверу.
type Request struct {
	// Поле Command может принимать четыре значения:
	// * "quit" - прощание с сервером (после этого сервер рвёт соединение);
	// * "add" - добавление пары
	Command string `json:"command"`

	// Если Command == "quit" поле Data пустое
	// В противном случае в поле Data должна лежать строка
	Data *json.RawMessage `json:"data"`

}

// Response -- ответ сервера клиенту.
type Response struct {
	// Поле Status может принимать три значения:
	// * "ok" - успешное выполнение команды "quit" , "vecs";
	// * "failed" - в процессе выполнения команды произошла ошибка;
	// * "result" - получено скалярное произведение
	Status string `json:"status"`

	// Если Status == "failed", то в поле Data находится сообщение об ошибке.
	Data *json.RawMessage `json:"data"`
}

//структура из ключа , значения и стороны (л или п)
type MapPeer struct {
	Side int `json:"side"`
	Key string `json:"key"`
	Value string `json:"value"`
}
