<h1 align="center"><img src="./public/brasilapi-logo-small.png"></h1>

# Consulta de CEP, CNPJ e ISPB Bancos

Este projeto é uma aplicação web que permite consultar informações sobre CEP, CNPJ e códigos bancários (ISPB) etc, usando a API da BrasilAPI. O backend é implementado em Go, e o frontend utiliza HTMX e Tailwind CSS para uma experiência de usuário moderna e responsiva.

## Tecnologias Utilizadas

- **Backend**: Go
- **Frontend**: HTMX, Tailwind CSS
- **API**: BrasilAPI (para consulta de CEP, CNPJ, ISPB Bancos, etc.)
## Funcionalidades

- **Consulta de CEP**: Permite buscar informações detalhadas sobre um CEP.
- **Consulta de CNPJ**: Permite buscar informações sobre uma empresa pelo CNPJ.
- **Consulta de ISPB Bancos**: Permite buscar informações sobre bancos utilizando o código ISPB.


## Descrição dos Componentes

- **main.go**: O arquivo principal para iniciar o servidor Go. Configura e executa o servidor HTTP.
- **internal/handlers/**: Contém manipuladores para lidar com as solicitações de API:
  - `cep_handler.go`: Manipulador para consultas de CEP.
  - `cnpj_handler.go`: Manipulador para consultas de CNPJ.
  - `bank_handler.go`: Manipulador para consultas de ISPB Bancos.
- **internal/models/**: Contém os modelos de dados para a aplicação:
  - `cep.go`: Modelo para os dados de CEP.
  - `cnpj.go`: Modelo para os dados de CNPJ.
  - `bank.go`: Modelo para os dados de ISPB Bancos.
- **internal/services/**: Contém os serviços que fazem as chamadas para a API da BrasilAPI e processam os dados:
  - `cep_service.go`: Serviço para buscar dados de CEP.
  - `cnpj_service.go`: Serviço para buscar dados de CNPJ.
  - `bank_service.go`: Serviço para buscar dados de ISPB Bancos.
- **web/templates/index.html**: Template HTML para a página inicial da aplicação.
- **go.mod**: Arquivo de módulo Go que gerencia as dependências do projeto.

Essa estrutura proporciona uma organização clara e modular do código, facilitando a manutenção e a expansão futura da aplicação.

## Como Executar o Projeto

1. **Clone o repositório:**

    ```bash
    git clone https://github.com/devluanpereira/ToolsDev.git
    cd ToolsDev
    ```

2. **Instale as dependências do Go:**

    ```bash
    go mod tidy
    ```

3. **Execute o servidor:**

    ```bash
    go run main.go
    ```

4. **Acesse a aplicação:**

    Abra seu navegador e vá para `http://localhost:8000` para ver a aplicação em funcionamento. Mais caso a porta esteja em uso mude para que esteja disponivel em `main.go`.


## Como Contribuir

1. Faça um fork deste repositório.
2. Crie uma branch para suas alterações (`git checkout -b minha-alteracao`).
3. Faça as alterações e commit (`git commit -am 'Adiciona minha alteração'`).
4. Envie para o repositório remoto (`git push origin minha-alteracao`).
5. Abra um Pull Request.

## Licença

Este projeto está licenciado sob a [MIT License](LICENSE).

## Contato

Se você tiver alguma dúvida, sinta-se à vontade para entrar em contato.

- **Nome:** Luan Pereira
- **Email:** luan23107@gmail.com
