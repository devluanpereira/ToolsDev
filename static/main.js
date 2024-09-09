const menuBtn = document.getElementById('menu-btn');
const mobileMenu = document.getElementById('mobile-menu');

menuBtn.addEventListener('click', () => {
    mobileMenu.classList.toggle('hidden');
});

// Mostrar e esconder o modal de CEP
const cepCard = document.getElementById('cep-card');
const cepModal = document.getElementById('cep-modal');
const closeCepModal = document.getElementById('close-cep-modal');

cepCard.addEventListener('click', () => {
    cepModal.classList.remove('hidden');
});

closeCepModal.addEventListener('click', () => {
    cepModal.classList.add('hidden');
});

// Mostrar e esconder o Modal de ISPB Bank
const bankCard = document.getElementById('bank-card');
const bankModal = document.getElementById('bank-modal');
const closeBankModal = document.getElementById('close-bank-modal');

bankCard.addEventListener('click', () => {
    bankModal.classList.remove('hidden');
});

closeBankModal.addEventListener('click', () => {
    bankModal.classList.add('hidden');
})


// Mostrar e esconder o modal de CNPJ
const cnpjCard = document.getElementById('cnpj-card');
const cnpjModal = document.getElementById('cnpj-modal');
const closeCnpjModal = document.getElementById('close-cnpj-modal');

cnpjCard.addEventListener('click', () => {
    cnpjModal.classList.remove('hidden');
});

closeCnpjModal.addEventListener('click', () => {
    cnpjModal.classList.add('hidden');
});

// Consulta de CEP
document.getElementById('cep-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    const cep = document.getElementById('cep').value;
    if (!cep) {
        alert("Por favor, digite um CEP válido.");
        return;
    }

    try {
        const response = await fetch(`/buscar-cep?cep=${cep}`);
        const data = await response.json();

        if (response.ok) {
            document.getElementById('cep-result').innerHTML = `
                        <p><strong>CEP:</strong> ${data.cep}</p>
                        <p><strong>Estado:</strong> ${data.state}</p>
                        <p><strong>Cidade:</strong> ${data.city}</p>
                        <p><strong>Bairro:</strong> ${data.neighborhood}</p>
                        <p><strong>Rua:</strong> ${data.street}</p>
                    `;
        } else {
            document.getElementById('cep-result').innerHTML = `<p class="text-red-500">Erro: ${data.message || 'CEP não encontrado.'}</p>`;
        }
    } catch (error) {
        document.getElementById('cep-result').innerHTML = `<p class="text-red-500">Erro ao buscar CEP: ${error.message}</p>`;
    }
});

// Consulta de CNPJ
document.getElementById('cnpj-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    const cnpj = document.getElementById('cnpj').value;
    if (!cnpj) {
        alert("Por favor, digite um CNPJ válido.");
        return;
    }

    try {
        const response = await fetch(`/buscar-cnpj?cnpj=${cnpj}`);
        const data = await response.json();

        if (response.ok) {
            document.getElementById('cnpj-result').innerHTML = `
                        <p><strong>CNPJ:</strong> ${data.cnpj}</p>
                        <p><strong>Razão Social:</strong> ${data.razao_social || 'Não disponível'}</p>
                        <p><strong>Nome Fantasia:</strong> ${data.nome_fantasia || 'Não disponível'}</p>
                        <p><strong>UF:</strong> ${data.uf || 'Não disponível'}</p>
                        <p><strong>Município:</strong> ${data.municipio || 'Não disponível'}</p>
                        <p><strong>Bairro:</strong> ${data.bairro || 'Não disponível'}</p>
                        <p><strong>Logradouro:</strong> ${data.logradouro || 'Não disponível'}</p>
                        <p><strong>Número:</strong> ${data.numero || 'Não disponível'}</p>
                        <p><strong>CEP:</strong> ${data.cep || 'Não disponível'}</p>
                    `;
        } else {
            document.getElementById('cnpj-result').innerHTML = `<p class="text-red-500">Erro: ${data.message || 'CNPJ não encontrado.'}</p>`;
        }
    } catch (error) {
        document.getElementById('cnpj-result').innerHTML = `<p class="text-red-500">Erro ao buscar CNPJ: ${error.message}</p>`;
    }
});

// Consulta de ISPB
document.getElementById('bank-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    const bank = document.getElementById('bank').value;
    if (!bank) {
        alert("Por favor digite um código válido");
        return;
    }
    try {
        const response = await fetch(`/buscar-code?code=${bank}`);
        const data = await response.json();

        if (response.ok) {
            document.getElementById('code-result').innerHTML = `
                        <p><strong>ISPB:</strong> ${data.ispb}</p>
                        <p><strong>Name:</strong> ${data.name}</p>
                        <p><strong>Nome Completo:</strong> ${data.fullname}</p>
                    `;
        } else {
            document.getElementById('code-result').innerHTML = `<p class="text-red-500">Erro: ${data.message || 'Código não encontrado.'}</p>`;
        }
    } catch (error) {
        document.getElementById('code-result').innerHTML = `<p class="text-red-500">Erro ao buscar ISPB: ${error.message}</p>`;
    }
});