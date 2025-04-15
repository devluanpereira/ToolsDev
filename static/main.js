// Menu Toggle
const menuBtn = document.getElementById('menu-btn');
const mobileMenu = document.getElementById('mobile-menu');
menuBtn.addEventListener('click', () => mobileMenu.classList.toggle('hidden'));

// Modal Toggle Helpers
function showModal(modal) {
    modal.classList.remove('hidden');
}

function hideModal(modal, clearFields = true) {
    modal.classList.add('hidden');
    if (clearFields) limparCampos();
}

// CEP Modal
const cepCard = document.getElementById('cep-card');
const cepModal = document.getElementById('cep-modal');
const closeCepModal = document.getElementById('close-cep-modal');

cepCard.addEventListener('click', () => showModal(cepModal));
closeCepModal.addEventListener('click', () => hideModal(cepModal));

// Bank Modal
const bankCard = document.getElementById('bank-card');
const bankModal = document.getElementById('bank-modal');
const closeBankModal = document.getElementById('close-bank-modal');

bankCard.addEventListener('click', () => showModal(bankModal));
closeBankModal.addEventListener('click', () => hideModal(bankModal));

// CNPJ Modal
const cnpjCard = document.getElementById('cnpj-card');
const cnpjModal = document.getElementById('cnpj-modal');
const closeCnpjModal = document.getElementById('close-cnpj-modal');

cnpjCard.addEventListener('click', () => showModal(cnpjModal));
closeCnpjModal.addEventListener('click', () => hideModal(cnpjModal));

// CEP Lookup
document.getElementById('cep-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    await performLookup('/buscar-cep', 'cep', 'cep-result', 'CEP');
});

// CNPJ Lookup
document.getElementById('cnpj-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    await performLookup('/buscar-cnpj', 'cnpj', 'cnpj-result', 'CNPJ');
});

// ISPB Lookup
document.getElementById('bank-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    await performLookup('/buscar-code', 'code', 'code-result', 'código');
});

// Lookup Helper Function
async function performLookup(apiUrl, inputId, resultId, label) {
    const input = document.getElementById(inputId).value;
    if (!input) {
        alert(`Por favor, digite um ${label} válido.`);
        return;
    }

    try {
        const response = await fetch(`${apiUrl}?${inputId}=${input}`);
        const data = await response.json();
        const resultContainer = document.getElementById(resultId);

        if (response.ok) {
            resultContainer.innerHTML = Object.entries(data)
                .map(([key, value]) => `<p><strong>${capitalize(key)}:</strong> ${value || 'Não disponível'}</p>`)
                .join('');
        } else {
            resultContainer.innerHTML = `<p class="text-red-500">Erro: ${data.message || `${label} não encontrado.`}</p>`;
        }
    } catch (error) {
        document.getElementById(resultId).innerHTML = `<p class="text-red-500">Erro ao buscar ${label}: ${error.message}</p>`;
    }
}

// Clear Fields
function limparCampos() {
    ['cnpj', 'cep', 'bank'].forEach(id => (document.getElementById(id).value = ''));
    ['cnpj-result', 'cep-result', 'code-result'].forEach(id => (document.getElementById(id).innerText = ''));
}

// Capitalize Helper
function capitalize(str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
}

// Search Functionality
document.getElementById('search-input').addEventListener('input', function () {
    const query = this.value.trim().toLowerCase();
    const cards = document.querySelectorAll('[data-title]');

    cards.forEach(card => {
        const title = card.getAttribute('data-title').toLowerCase();
        card.style.display = title.includes(query) ? 'block' : 'none';
    });

    if (!query) cards.forEach(card => (card.style.display = 'block'));
});
