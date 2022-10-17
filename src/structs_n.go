package chat_bot_connector

// internal gRPC structure

/*
Метод Find
аргументы:
	query - строка - идентификатор запроса
	desired - структура - то что ищем
	essence - строка - имя сущности
	property - строка - имя свойства
	conditions - массив структур - список условий
		condition - структура - условие
			essence - строка - имя сущности
			property - строка - имя свойства
			value - строка - значение свойства

возврат
	result - строка - OK, Error
	error - строка - описание ошибки
	founds - массив структур - список найденных
		found - структура - найденное
			essence - строка - имя сущности
			property - строка - имя свойства
			value - строка - значение свойства

*/

type FindCondition struct {
	Essence  string
	Property string
	Value    string
}

type FindFound struct {
	Essence  string
	Property string
	Value    string
}

// структуры внутренней обработки очереди
type Command struct {
	Cmd      string // команда
	Sentence string // предложение для разбора
	ID       string // идентификатор в очереди обработки
}

type CommandAnswer struct {
	Cmd      string   // команда
	Sentence string   // входная фраза
	WordLink []string // обработанное предложение
	ID       string   // идентификатор в очереди обработки
	Result   string
	Error    string
}

type Settings struct {
	AddrClient       string // выносим в переменную
	PortClient       int    // выносим в переменную
	PortServer       int    // выносим в переменную
	PathFileStorage  string // в переменной окружения
	DBEnv            string // реквизиты доступа к базе данных
	SettingsFileName string // в переменной окружения
	// мы хотим конфигурировать микросервис из основного сервиса?
//	Files *SettingsFile
}
