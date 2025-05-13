# Weather-By-CEP-Otel

Esta aplicação permite consultar a temperatura atual em graus Celsius, Fahrenheit, Kelvin, pesquisando por CEP em um endpoint **REST**. A aplicação está dividida em dois serviços, o A (porta :8181) que recebe o input e trata o CEP e o serviço B (porta :8282) que orquestra as consultas de localização e temperatura, com tracing distribuído entre os seviços e as métricas de tempo da execução dos serviços e das chamadas externas executadas.

## Executando a aplicação em ambiente local
1. Certifique-se de ter o Docker instalado.
2. Suba os containers necessários executando o comando:
    ```bash
    docker-compose up --build
    ```
3. Aguarde até que a mensagem de que as aplicações estejam rodando nas portas :8181 e :8282 seja exibida nos logs.
4. Pronto! O serviço esta disponível no ambiente local. Pode ser consumido usando o modelo disponível em `api/temp_local.http` (ajustar o CEP).

## Tracing distribuído
Ao executar a aplicação, também subimos a implementação do tracing distribuído com OTEL e Zipkin. Essa implementação nos possibilita acompanhar a execução de ambos os serviços, assim como suas chamadas externas.

* Os traces podem ser visualizados pelo painel do Zipkin, acessando `http://localhost:9412/`.

Abertura dos traces:
1. service-a-request - Execução do serviço A (input)
2. fetch-service-b - Chamada do serviço B pelo serviço A
3. service-b-request - Execução do serviço B (orchestrator)
4. fetch-location-service - Chamada do serviço ViaCEP pelo serviço B
5. fetch-weather-service - Chamada do serviço WeatherAPI pelo serviço B