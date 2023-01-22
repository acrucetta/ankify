const form = document.querySelector("#form") as HTMLFormElement;
const input = document.querySelector("#input") as HTMLInputElement;
const output = document.querySelector("#output") as HTMLElement;
const loading = document.querySelector("#loading") as HTMLElement;

form.addEventListener("submit", async (e) => {
    e.preventDefault();
    form.classList.add("loading");
    const text = input.value;
    try {
        const response = await fetch("http://localhost:8080/getquestions", {
            method: "POST",
            body: JSON.stringify({ text }),
            headers: { "Content-Type": "application/json" }
        });
        const result = await response.json();
        output.innerHTML = result.output;
    } catch (err) {
        console.error(err);
    } finally {
        form.classList.remove("loading");
    }
});
