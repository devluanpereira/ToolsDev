document.getElementById("pagamento-form").addEventListener("submit", async function (e) {
    e.preventDefault();

    const form = e.target;
    const formData = new FormData(form);
    const params = new URLSearchParams();

    for (const pair of formData.entries()) {
      params.append(pair[0], pair[1]);
    }

    try {
      const response = await fetch("/criar-pagamento", { // Requisição para o seu servidor
        method: "POST",
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: params.toString(),
        credentials: "include"
      });

      if (!response.ok) {
        const errorText = await response.text();
        alert(`Erro ao criar pagamento: ${response.status} - ${errorText}`);
        console.error(`Erro ao criar pagamento: ${response.status} - ${errorText}`);
        return;
      }

      const data = await response.json(); // Espera uma resposta JSON do seu servidor com a init_point

      if (data && data.init_point) {
        window.location.href = data.init_point; // Redireciona para o Mercado Pago
      } else {
        alert("Resposta inválida do servidor.");
        console.error("Resposta inválida do servidor:", data);
      }

    } catch (err) {
      alert("Erro ao processar pagamento.");
      console.error(err);
    }
});
