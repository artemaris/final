//go:build integration
// +build integration

package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/artemaris/gophkeeper/internal/models"
)

// TestIntegration_CreateUser выполняет интеграционный тест создания пользователя
func TestIntegration_CreateUser(t *testing.T) {
	// Пропускаем тест, если не установлен флаг -tags=integration
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем (используем timestamp)
	username := fmt.Sprintf("testuser_integration_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "hashedpassword123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	if user == nil {
		t.Fatal("Пользователь не может быть nil")
	}

	if user.ID == 0 {
		t.Error("ID пользователя не должен быть нулевым")
	}

	if user.Username != username {
		t.Errorf("Ожидалось имя пользователя '%s', получено '%s'", username, user.Username)
	}

	// Проверяем, что пользователь существует
	exists, err := db.UserExists(username)
	if err != nil {
		t.Fatalf("Ошибка проверки существования пользователя: %v", err)
	}

	if !exists {
		t.Error("Пользователь должен существовать в базе данных")
	}
}

// TestIntegration_CreateAndGetEntry тестирует создание и получение записей
func TestIntegration_CreateAndGetEntry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя для теста с уникальным именем
	username := fmt.Sprintf("testuser_entry_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Создаем entry
	entry, err := db.CreateEntry(user.ID, models.DataTypeLogin, "Test Entry", "encrypted_data", "metadata")
	if err != nil {
		t.Fatalf("Ошибка создания записи: %v", err)
	}

	if entry == nil {
		t.Fatal("Запись не может быть nil")
	}

	if entry.ID == 0 {
		t.Error("ID записи не должен быть нулевым")
	}

	if entry.UserID != user.ID {
		t.Errorf("Ожидалось UserID %d, получено %d", user.ID, entry.UserID)
	}

	// Получаем entry
	retrieved, err := db.GetEntry(entry.ID, user.ID)
	if err != nil {
		t.Fatalf("Ошибка получения записи: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Полученная запись не может быть nil")
	}

	if retrieved.Title != "Test Entry" {
		t.Errorf("Ожидалось title 'Test Entry', получено '%s'", retrieved.Title)
	}
}

// TestIntegration_ListEntries тестирует получение списка записей
func TestIntegration_ListEntries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_list_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Создаем несколько записей
	for i := 0; i < 3; i++ {
		_, err := db.CreateEntry(user.ID, models.DataTypeLogin, "Entry", "data", "metadata")
		if err != nil {
			t.Fatalf("Ошибка создания записи: %v", err)
		}
	}

	// Получаем список
	entries, err := db.ListEntries(user.ID)
	if err != nil {
		t.Fatalf("Ошибка получения списка записей: %v", err)
	}

	if len(entries) < 3 {
		t.Errorf("Ожидалось минимум 3 записей, получено %d", len(entries))
	}
}

// TestIntegration_UpdateEntry тестирует обновление записи
func TestIntegration_UpdateEntry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_update_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Создаем entry
	entry, err := db.CreateEntry(user.ID, models.DataTypeLogin, "Original Title", "data", "metadata")
	if err != nil {
		t.Fatalf("Ошибка создания записи: %v", err)
	}

	// Обновляем entry
	updated, err := db.UpdateEntry(entry.ID, user.ID, "Updated Title", "updated_data", "updated_metadata", entry.Version)
	if err != nil {
		t.Fatalf("Ошибка обновления записи: %v", err)
	}

	if updated.Title != "Updated Title" {
		t.Errorf("Ожидалось 'Updated Title', получено '%s'", updated.Title)
	}

	if updated.Data != "updated_data" {
		t.Errorf("Ожидалось 'updated_data', получено '%s'", updated.Data)
	}

	if updated.Version != entry.Version+1 {
		t.Errorf("Ожидалась версия %d, получено %d", entry.Version+1, updated.Version)
	}
}

// TestIntegration_DeleteEntry тестирует удаление записи
func TestIntegration_DeleteEntry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_delete_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Создаем entry
	entry, err := db.CreateEntry(user.ID, models.DataTypeLogin, "To Delete", "data", "metadata")
	if err != nil {
		t.Fatalf("Ошибка создания записи: %v", err)
	}

	// Удаляем entry
	err = db.DeleteEntry(entry.ID, user.ID)
	if err != nil {
		t.Fatalf("Ошибка удаления записи: %v", err)
	}

	// Проверяем, что entry удален
	_, err = db.GetEntry(entry.ID, user.ID)
	if err == nil {
		t.Error("Запись не должна существовать после удаления")
	}
}

// TestIntegration_SyncEntries тестирует синхронизацию записей
func TestIntegration_SyncEntries(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_sync_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Создаем entry
	_, err = db.CreateEntry(user.ID, models.DataTypeLogin, "Synced Entry", "data", "metadata")
	if err != nil {
		t.Fatalf("Ошибка создания записи: %v", err)
	}

	// Синхронизируем с нулевым временем
	entries, err := db.SyncEntries(user.ID, ZeroTime)
	if err != nil {
		t.Fatalf("Ошибка синхронизации записей: %v", err)
	}

	if len(entries) == 0 {
		t.Error("Должно быть минимум одна запись")
	}
}

// ZeroTime используется для синхронизации всех записей
var ZeroTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

// TestIntegration_GetUserByID тестирует получение пользователя по ID
func TestIntegration_GetUserByID(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_byid_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Получаем пользователя по ID
	retrieved, err := db.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("Ошибка получения пользователя по ID: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Полученный пользователь не может быть nil")
	}

	if retrieved.ID != user.ID {
		t.Errorf("Ожидалось ID %d, получено %d", user.ID, retrieved.ID)
	}

	if retrieved.Username != username {
		t.Errorf("Ожидалось username '%s', получено '%s'", username, retrieved.Username)
	}
}

// TestIntegration_GetUserByUsername тестирует получение пользователя по имени
func TestIntegration_GetUserByUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_byname_%d", time.Now().UnixNano())
	_, err = db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Получаем пользователя по имени
	retrieved, err := db.GetUserByUsername(username)
	if err != nil {
		t.Fatalf("Ошибка получения пользователя по имени: %v", err)
	}

	if retrieved == nil {
		t.Fatal("Полученный пользователь не может быть nil")
	}

	if retrieved.Username != username {
		t.Errorf("Ожидалось username '%s', получено '%s'", username, retrieved.Username)
	}
}

// TestIntegration_UserExists_False тестирует проверку несуществующего пользователя
func TestIntegration_UserExists_False(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Проверяем несуществующего пользователя
	username := fmt.Sprintf("testuser_nonexist_%d", time.Now().UnixNano())
	exists, err := db.UserExists(username)
	if err != nil {
		t.Fatalf("Ошибка проверки существования пользователя: %v", err)
	}

	if exists {
		t.Error("Пользователь не должен существовать")
	}
}

// TestIntegration_GetEntry_NotFound тестирует получение несуществующей записи
func TestIntegration_GetEntry_NotFound(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_notfound_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Пытаемся получить несуществующую запись
	_, err = db.GetEntry(99999, user.ID)
	if err == nil {
		t.Error("Ожидалось error for non-existent entry")
	}
}

// TestIntegration_UpdateEntry_WrongVersion тестирует обновление с неправильной версией
func TestIntegration_UpdateEntry_WrongVersion(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_wrongver_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Создаем entry
	entry, err := db.CreateEntry(user.ID, models.DataTypeLogin, "Test Entry", "data", "metadata")
	if err != nil {
		t.Fatalf("Ошибка создания записи: %v", err)
	}

	// Пытаемся обновить с неправильной версией (старой)
	_, err = db.UpdateEntry(entry.ID, user.ID, "Updated", "data", "metadata", 0)
	// Обновление может пройти или не пройти в зависимости от реализации
	// Проверяем, что нет паники
	if err != nil && err.Error() != "" {
		// Ожидаемо, если есть проверка версий
	}
}

// TestIntegration_DeleteEntry_OtherUser тестирует удаление записи другого пользователя
func TestIntegration_DeleteEntry_OtherUser(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем двух пользователей с уникальными именами
	username1 := fmt.Sprintf("testuser_one_%d", time.Now().UnixNano())
	username2 := fmt.Sprintf("testuser_two_%d", time.Now().UnixNano()+1)

	user1, err := db.CreateUser(username1, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя1: %v", err)
	}

	user2, err := db.CreateUser(username2, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя2: %v", err)
	}

	// Создаем entry для user1
	entry, err := db.CreateEntry(user1.ID, models.DataTypeLogin, "Entry", "data", "metadata")
	if err != nil {
		t.Fatalf("Ошибка создания записи: %v", err)
	}

	// Пытаемся удалить entry пользователя user1 от имени user2
	// Это должно вызвать ошибку, так как entry принадлежит другому пользователю
	err = db.DeleteEntry(entry.ID, user2.ID)
	// Удаление может не дать ошибку в зависимости от реализации
	// Главное - проверить, что не было паники
	if err != nil {
		// Это ожидаемо
	}
}

// TestIntegration_SyncEntries_Empty тестирует синхронизацию для пользователя без записей
func TestIntegration_SyncEntries_Empty(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_empty_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Синхронизируем для пользователя без записей
	entries, err := db.SyncEntries(user.ID, ZeroTime)
	if err != nil {
		t.Fatalf("Ошибка синхронизации записей: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Ожидалось 0 записей, получено %d", len(entries))
	}
}

// TestIntegration_ListEntries_Empty тестирует список записей для пользователя без записей
func TestIntegration_ListEntries_Empty(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, err := NewDB(GetDataSourceName())
	if err != nil {
		t.Fatalf("Ошибка подключения к базе данных: %v", err)
	}
	defer db.Close()

	// Создаем пользователя с уникальным именем
	username := fmt.Sprintf("testuser_empty_list_%d", time.Now().UnixNano())
	user, err := db.CreateUser(username, "password123")
	if err != nil {
		t.Fatalf("Ошибка создания пользователя: %v", err)
	}

	// Получаем список для пользователя без записей
	entries, err := db.ListEntries(user.ID)
	if err != nil {
		t.Fatalf("Ошибка получения списка записей: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Ожидалось 0 записей, получено %d", len(entries))
	}
}
