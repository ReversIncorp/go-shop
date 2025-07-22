### Введение
Для работы dendrite необходим приватный ключ и конфиг. 

### Генерация стандартного конфига
```bash
CONFIG_OUTPUT=$(docker run --rm --entrypoint=/usr/bin/generate-config \
  -v $(pwd)/dendrite:/etc/dendrite \
  matrixdotorg/dendrite-monolith:latest \
  -dir /etc/dendrite)
echo "$CONFIG_OUTPUT" > dendrite_config.yaml
```

### Генерация приватного ключа
```bash
mkdir -p dendrite
docker run --rm --entrypoint=/usr/bin/generate-keys \
  -v $(pwd)/dendrite:/etc/dendrite \
  matrixdotorg/dendrite-monolith:latest \
  --private-key /etc/dendrite/matrix_key.pem
cp dendrite/matrix_key.pem ./dendrite_matrix_key.pem
rm -rf dendrite
```