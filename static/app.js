(function () {
  function byId(id) {
    return document.getElementById(id);
  }

  function parseEnvelope(payload) {
    if (Array.isArray(payload)) {
      return payload;
    }

    if (payload && Array.isArray(payload.data)) {
      return payload.data;
    }

    return [];
  }

  function transactionAmount(transaction) {
    if (typeof transaction.amount === "string" && transaction.amount.trim() !== "") {
      return transaction.amount;
    }

    const cents = Number(
      transaction.amount_in_cents ??
      transaction.AmountInCents ??
      transaction.amountInCents ??
      transaction.amount_cents ??
      0
    );

    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
    }).format(cents / 100);
  }

  function transactionName(transaction, index) {
    const value = transaction.name ?? transaction.Name ?? transaction.title ?? transaction.label;
    if (typeof value === "string" && value.trim() !== "") {
      return value.trim();
    }

    return "";
  }

  function transactionType(transaction) {
    return String(transaction.type ?? transaction.Type ?? "DEBIT").toUpperCase();
  }

  function transactionSlug(transaction, index) {
    return transaction.slug ?? transaction.Slug ?? transaction.id ?? transaction.ID ?? "txn-" + String(index + 1);
  }

  function renderList(transactions) {
    const list = byId("transaction-list");
    const emptyState = byId("empty-state");
    const count = byId("transaction-count");
    const heroCount = byId("hero-transaction-count");

    count.textContent = String(transactions.length);
    if (heroCount) {
      heroCount.textContent = String(transactions.length);
    }

    if (!transactions.length) {
      list.hidden = true;
      list.innerHTML = "";
      emptyState.hidden = false;
      return;
    }

    emptyState.hidden = true;
    list.hidden = false;
    list.innerHTML = transactions.map(function (transaction, index) {
      const type = transactionType(transaction);
      const typeClass = type.toLowerCase();
      const name = escapeHTML(transactionName(transaction, index));
      const amount = escapeHTML(transactionAmount(transaction));
      const slug = escapeHTML(String(transactionSlug(transaction, index)));

      return [
        '<li class="transaction-item">',
        '  <div class="transaction-copy">',
        '    <p class="transaction-name">' + name + "</p>",
        '    <p class="transaction-amount">' + amount + "</p>",
        '    <p class="transaction-slug">Reference #' + slug + "</p>",
        "  </div>",
        '  <span class="transaction-type transaction-type-' + typeClass + '">' + escapeHTML(type) + "</span>",
        "</li>",
      ].join("");
    }).join("");
  }

  function escapeHTML(value) {
    return value
      .replaceAll("&", "&amp;")
      .replaceAll("<", "&lt;")
      .replaceAll(">", "&gt;")
      .replaceAll('"', "&quot;")
      .replaceAll("'", "&#39;");
  }

  function setStatus(message, state) {
    const banner = byId("status-banner");
    banner.hidden = false;
    banner.className = "status-banner" + (state ? " is-" + state : "");
    banner.textContent = message;
  }

  function clearStatus() {
    const banner = byId("status-banner");
    banner.hidden = true;
    banner.textContent = "";
    banner.className = "status-banner";
  }

  function setFormError(message) {
    const error = byId("form-error");
    if (!message) {
      error.hidden = true;
      error.textContent = "";
      return;
    }

    error.hidden = false;
    error.textContent = message;
  }

  function extractError(payload, fallback) {
    if (payload && payload.error && payload.error.message) {
      return payload.error.message;
    }

    if (payload && typeof payload.message === "string" && payload.message.trim() !== "") {
      return payload.message;
    }

    return fallback;
  }

  async function loadTransactions(listEndpoint) {
    setStatus("Loading transactions...", "loading");

    const response = await fetch(listEndpoint, {
      headers: {
        Accept: "application/json",
      },
    });

    const payload = await response.json().catch(function () {
      return null;
    });

    if (!response.ok) {
      throw new Error(extractError(payload, "Failed to load transactions."));
    }

    renderList(parseEnvelope(payload));
    clearStatus();
  }

  async function createTransaction(createEndpoint, payload) {
    const response = await fetch(createEndpoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Accept: "application/json",
      },
      body: JSON.stringify(payload),
    });

    const data = await response.json().catch(function () {
      return null;
    });

    if (!response.ok) {
      throw new Error(extractError(data, "Failed to save transaction."));
    }
  }

  async function boot() {
    const root = byId("transactions-app");
    if (!root) {
      return;
    }

    const listEndpoint = root.dataset.listEndpoint;
    const createEndpoint = root.dataset.createEndpoint;
    const form = byId("transaction-form");

    try {
      await loadTransactions(listEndpoint);
    } catch (error) {
      setStatus(error.message, "error");
    }

    form.addEventListener("submit", async function (event) {
      event.preventDefault();
      setFormError("");

      const formData = new FormData(form);
      const payload = {
        name: String(formData.get("name") || "").trim(),
        amount: String(formData.get("amount") || "").trim(),
        type: String(formData.get("type") || "").trim(),
      };

      try {
        setStatus("Saving transaction...", "loading");
        await createTransaction(createEndpoint, payload);
        form.reset();
        form.elements.type.value = "DEBIT";
        await loadTransactions(listEndpoint);
      } catch (error) {
        setFormError(error.message);
        setStatus(error.message, "error");
      }
    });
  }

  window.addEventListener("DOMContentLoaded", boot);
})();
