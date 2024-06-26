package ring

import (
	"fmt"
	"sync"
)

// Структура кольцевого буфера
type IntBuffer struct {
	m     sync.Mutex // мьютекс для потокобезопасного доступа к буферу, т.к. оепрации чтения/записи могут производиться разными горутинами
	data  []*int     // само хранилище данных
	size  int        // размер буфера
	used  int        // использованно памяти
	read  int        // указатель на ячейку для чтения
	write int        // указатель на ячейку для записи
}

// Конструктор кольцевого буфера, стартовая позиция равна началу реального хранилища
func NewIntBuffer(size int) *IntBuffer {
	return &IntBuffer{sync.Mutex{}, make([]*int, size), size, 0, 0, 0}
}

// Size - метод, возвращающий размер кольцевого массива
func (r *IntBuffer) Size() int {
	r.m.Lock()
	defer r.m.Unlock()
	return r.size
}

// Used - метод, возвращающий количество занятых ячеек
func (r *IntBuffer) Used() int {
	r.m.Lock()
	defer r.m.Unlock()
	return r.used
}

// IsEmpty - Метод, проверяющий кольцо на пустоту
func (r *IntBuffer) IsEmpty() bool {
	r.m.Lock()
	defer r.m.Unlock()
	return r.used == 0
}

// IsFull - Проверка на полноту, достигнут ли конец буфера
func (r *IntBuffer) IsFull() bool {
	r.m.Lock()
	defer r.m.Unlock()
	return r.used == r.size
}

// Read - метод чтения элемента
func (r *IntBuffer) Read() (int, error) {
	r.m.Lock()
	defer r.m.Unlock()
	if r.used == 0 { // Проверка на пустоту
		return 0, fmt.Errorf("буфер пустой")
	}
	el := *r.data[r.read]
	r.data[r.read] = nil           // очищаем ячейку
	r.read = (r.read + 1) % r.size // переводим указатель чтения на следующую ячейку буфера
	r.used--                       // уменьшаем кол-во использованных ячеек
	return el, nil
}

// Write - запись нового элемента в буфер
func (r *IntBuffer) Write(v int) error {
	r.m.Lock()
	defer r.m.Unlock()
	if r.used == r.size { // Проверка на заполненность
		return fmt.Errorf("буфер полон")
	}
	r.data[r.write] = &v             //записываем данные в ячейку
	r.write = (r.write + 1) % r.size // перемещаем указатель записи на следующую ячейку буфера
	r.used++                         // увеличиваем кол-во использованных ячеек буфера
	return nil
}

// Output - вывод всех элементов буфера и очистка данных
func (r *IntBuffer) Output() []int {
	r.m.Lock()
	defer r.m.Unlock()
	res := make([]int, 0)
	if r.used == 0 {
		return res
	}
	for { // перебираем хранилище от r.read до r.write
		res = append(res, *r.data[r.read]) // забираем значение
		r.data[r.read] = nil               // очищаем ячейку
		r.read = (r.read + 1) % r.size     // переходим  к следующей
		if r.read == r.write {             // если сравнялись с указателяем записи, выходим
			break
		}
	}
	r.used = 0 // обнуляем кол-во использованных ячеек буфера
	return res
}

func (r *IntBuffer) String() string {
	r.m.Lock()
	defer r.m.Unlock()
	res := "[\n"
	res += fmt.Sprintf(" Размер: %d\n", r.size)
	res += fmt.Sprintf(" Использовано: %d\n", r.used)
	res += fmt.Sprintf(" Ячейка для чтения: %d\n", r.read)
	res += fmt.Sprintf(" Ячейка для записи: %d\n", r.write)
	res += " Содержимое буфера:\n"
	for i, el := range r.data {
		if el != nil {
			res += fmt.Sprintf("  [%d: %d]\n", i, *el)
		} else {
			res += fmt.Sprintf("  [%d: пусто]\n", i)
		}
	}
	res += "]"
	return res
}
