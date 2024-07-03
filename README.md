# Weather Notification System
Weather Notification System é um sistema de notificação de temperatura, clima e tempo para determinadas regiões do Brasil. O sistema funciona de maneira genérica e permite ser 
integrado com diversos sistemas que queiram notificar seus usuários sobre o tempo.

Os usuários serão notificados com informações sobre a previsão do tempo para os próximos 4 dias e para os usuários que moram em regiões litorâneas também receberão informações sobre velocidade de vento, altura de ondas e etc.

A origem das informações são fonte do Centro de Previsão de Tempo e Estudos
Climáticos (CPTEC), para mais informações acesse o [site](https://www.cptec.inpe.br/sp/sao-paulo).

## Solução técnica
A solução é composta por 3 principais aplicações, **API's, Workers e Websockets (clients/servers)** com as tecnologias Golang, Postgres e RabbitMQ.

- **API:** aplicação para registrar usuários na plataforma e enviar notificações para um grupo de usuários
- **Worker:** aplicação para notificar os usuários em background
- **Webskcoket:** client/server que mantêm uma comunicação com os usuários conectados e fazendo a notificação para  cada um deles

Para garantir **resiliência** a aplicação trabalha com eventos que são processados em background, alêm disso as integrações com a CPTEC é garantida através de retry com backoff.
Para garantir a **escalabilidade** as aplicações são independentes uma das outras e podem rodar com mais de uma instância ao mesmo tempo.
A aplicação está preparada para rodar em container com apenas uma imagem Docker que inicia a aplicação de acordo com o entrypoint recebido (api, worker, websocket)

![image](https://drive.google.com/uc?export=view&id=1llDQ2xr3QLoubmjJC3B848W7KjpcLTpY)


## Use cases
Os principais use cases da aplicaçãos são:
- **Registar usuário**: Recebe informações como nome, sobrenome, e localidade. Faz a busca na API da CPTEC para verificar se a localização é suportada e caso seja salva o usuário na plataforma.
- **Unsubscribe usuário**: Remove um determinado usuário de receber notificações futuras ou já agendadas.
- **Notificar usuários**: Envia notificação sobre o tempo para uma lista de usuários. Essas notifiações serão enfileiradas e processadas uma a uma em background 

## Rodar local
  ```sh
make build
make lint
make deps/start
make run/api
make run/worker-web
make run/websocket-c
```

### author
[Guilherme Lyra](https://github.com/lyracampos "Guilherme Lyra")

[Linkedin](https://www.linkedin.com/in/guilherme-lyra-campos/)
