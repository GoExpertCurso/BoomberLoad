# BoomberLoad

## Introdução

Este projeto é uma aplicação cli que realiza testes de carga a serviços web com base nos parâmetros inseridos pelo usuário.

A aplicação funciona através de três parâmetros fornecidos pelo usuário.
  - URL: url do serviço web a ser realizado o teste.
  - requests: Número de requisições que serão realizadas para o serviço.
  - concurrency: Número de chamadas simultanêas que o aplicação devera fazer.

## Funcionalidade

### Stress Test
  * Objetivo: Realizar testes de carga em um serviço web
  * Parâmetro: 
    * **-url**: Endereço do serviço web.
    * **-requests**: Número de requisições desejadas.
    * **-concurrency**: Número de chamadas simultâneas.


## Execução
1. Tenha o docker instalado.
2. Execute o seguinte comando para fazer a construção da imagem:
    ```Bash
    docker build -t goexpert/boomberload .
    ```
3. Execute o seguinte comando para executar a ferramente de load test: 
   ```Bash
   docker run goexpert/boomberload -url=https://www.youtube.com/ -requests=1000 -concurrency=10
   ```
