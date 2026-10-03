#!/bin/bash
echo "========================================"
echo "  Тестирование REST API Celebrities"
echo "========================================"
echo

echo "[1] GET /Celebrities/All — вся коллекция"
echo "----------------------------------------"
curl -s -X GET http://localhost:3000/Celebrities/All
echo -e "\n"

echo "[2] GET /Celebrities/1 — элемент с id=1"
echo "----------------------------------------"
curl -s -X GET http://localhost:3000/Celebrities/1
echo -e "\n"

echo "[3] GET /Celebrities/999 — несуществующий id (ожидаем 404)"
echo "----------------------------------------"
curl -s -o /dev/null -w "HTTP Status: %{http_code}" -X GET http://localhost:3000/Celebrities/999
echo -e "\n"

echo "[4] POST /Celebrities — добавить новый элемент (ожидаем 201)"
echo "----------------------------------------"
curl -s -o /dev/null -w "HTTP Status: %{http_code}" -X POST http://localhost:3000/Celebrities \
  -H "Content-Type: application/json" \
  -d '{"id":10,"fullName":"Тест Тестов","nationality":"Тестовый","reqPhotoPath":"/photos/test.jpg"}'
echo -e "\n"

echo "[5] POST /Celebrities — дубликат id=10 (ожидаем 409)"
echo "----------------------------------------"
curl -s -o /dev/null -w "HTTP Status: %{http_code}" -X POST http://localhost:3000/Celebrities \
  -H "Content-Type: application/json" \
  -d '{"id":10,"fullName":"Дубликат","nationality":"Нет","reqPhotoPath":"/photos/dup.jpg"}'
echo -e "\n"

echo "[6] GET /Celebrities/All — проверяем что id=10 добавлен"
echo "----------------------------------------"
curl -s -X GET http://localhost:3000/Celebrities/All
echo -e "\n"

echo "[7] PUT /Celebrities/10 — обновить элемент с id=10 (ожидаем 200)"
echo "----------------------------------------"
curl -s -o /dev/null -w "HTTP Status: %{http_code}" -X PUT http://localhost:3000/Celebrities/10 \
  -H "Content-Type: application/json" \
  -d '{"fullName":"Обновлённый Тест","nationality":"Обновлённый","reqPhotoPath":"/photos/updated.jpg"}'
echo -e "\n"

echo "[8] PUT /Celebrities/999 — несуществующий id (ожидаем 404)"
echo "----------------------------------------"
curl -s -o /dev/null -w "HTTP Status: %{http_code}" -X PUT http://localhost:3000/Celebrities/999 \
  -H "Content-Type: application/json" \
  -d '{"fullName":"Нет","nationality":"Нет","reqPhotoPath":""}'
echo -e "\n"

echo "[9] GET /Celebrities/10 — проверяем обновление"
echo "----------------------------------------"
curl -s -X GET http://localhost:3000/Celebrities/10
echo -e "\n"

echo "[10] DELETE /Celebrities/10 — удалить id=10 (ожидаем 200)"
echo "----------------------------------------"
curl -s -o /dev/null -w "HTTP Status: %{http_code}" -X DELETE http://localhost:3000/Celebrities/10
echo -e "\n"

echo "[11] DELETE /Celebrities/999 — несуществующий id (ожидаем 404)"
echo "----------------------------------------"
curl -s -o /dev/null -w "HTTP Status: %{http_code}" -X DELETE http://localhost:3000/Celebrities/999
echo -e "\n"

echo "[12] GET /Celebrities/All — финальное состояние коллекции"
echo "----------------------------------------"
curl -s -X GET http://localhost:3000/Celebrities/All
echo -e "\n"

echo "========================================"
echo "  Тестирование завершено"
echo "========================================"
