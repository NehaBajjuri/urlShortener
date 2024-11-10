document.getElementById("shortenForm").addEventListener("submit", function (e) {
    e.preventDefault();

    const urlInput = document.getElementById("urlInput").value;

    fetch("/shorten", {
        method: "POST",
        headers: {
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: `url=${encodeURIComponent(urlInput)}`,
    })
    .then(response => response.text())
    .then(html => {
        document.getElementById("result").innerHTML = html;
    })
    .catch(error => console.error("Error:", error));
});
