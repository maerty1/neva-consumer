## Consumer

Репозиторий: https://github.com/urus-neva/consumer/tree/main

Сервисы запускаются в Docker через docker compose.

Запуск каждого сервиса происходит по единому принципу с помощью workflow в GitHub.

Все workflow расположены https://github.com/urus-neva/consumer/tree/main/.github/workflows

&nbsp;

# Логика работы

На примере wokflow https://github.com/urus-neva/consumer/blob/main/.github/workflows/bff_service.yml

**Блок**

```yaml
  deploy:
    runs-on: ubuntu-22.04
```

Запустить выполнение workflow в контейнере с ОС ubuntu 22.04.

&nbsp;

**Блок**

```yaml
on:
  push:
    paths:
      - 'bff_service/**'
    branches:
      - main
      - dev
```

Если в https://github.com/urus-neva/consumer/tree/main/bff_service будет сделан коммит, то сработает workflow.

&nbsp;

**Блок**

```yaml
      - name: Replace variables in docker-compose.yml
        uses: danielr1996/envsubst-action@1.0.0
        env:
            JWT_SECRET: ${{ secrets.JWT_SECRET }}
        if: ${{ github.ref == 'refs/heads/main' }}
        with:
          input: bff_service/docker/docker-compose.template.yml
          output: bff_service/docker/docker-compose.yml


      - name: Replace variables in dev-docker-compose.yml
        uses: danielr1996/envsubst-action@1.0.0
        env:
            JWT_SECRET: ${{ secrets.JWT_SECRET }}
        if: ${{ github.ref == 'refs/heads/dev' }}
        with:
          input: bff_service/docker/dev-docker-compose.template.yml
          output: bff_service/docker/dev-docker-compose.yml
```

В этом блоке в работу включается модуль `danielr1996/envsubst-action@1.0.0`, который формирует из шаблона `bff_service/docker/docker-compose.template.yml` рабочий `bff_service/docker/docker-compose.yml`.

Этот модуль нужен для того, чтобы GitHub мог вставлять содержимое переменных (https://github.com/urus-neva/consumer/settings/secrets/actions) в файл `docker-compose.yml`.

В зависимости от ветки срабатывает тот или иной блок.

&nbsp;

**Блок**

```yaml
      - name: Install SSH keys
        run: |
          install -m 600 -D /dev/null ~/.ssh/id_rsa
          echo "${{ secrets.SSH_PRIVATE_KEY }}" > ~/.ssh/id_rsa
          ssh-keyscan -H ${{ secrets.SSH_HOST }} > ~/.ssh/known_hosts

      - name: Create Docker context
        run: |
          docker context create remote \
            --docker host=ssh://${{ secrets.SSH_USER }}@${{ secrets.SSH_HOST }} \
            --description "Remote Docker context for ${{ secrets.SSH_HOST }}"
```

В этом блоке извлекается закрытый ключ из переменной `SSH_PRIVATE_KEY`,  по которому происходит подключение к хосту по SSH, который указан в переменной `SSH_HOST`. Далее, docker создаёт контекст с именем remote. Через этот контекст в дальнейшем выполняется запуск `docker-compose.yml` на удаленном хосте.

&nbsp;

**Блок**

```yaml
      - name: Docker compose
        run: |
          docker --context remote compose -f bff_service/docker/docker-compose.yml up -d --build
        if: ${{ github.ref == 'refs/heads/main' }}

      - name: Docker compose (dev)
        run: |
          docker --context remote compose -f bff_service/docker/dev-docker-compose.yml up -d --build
        if: ${{ github.ref == 'refs/heads/dev' }}
```

В этом блоке docker запускает `docker-compose.yml` на удаленном хосте через контекст remote. В зависимости от ветки выполняется тот или иной блок.

&nbsp;

**Блок**

```yaml
      - name: Cleanup
        run: rm -rf ~/.ssh
```

Удаляет папку `~/.ssh` в контейнере Ubuntu 22.04, в котором выполняется wokrflow.

&nbsp;
