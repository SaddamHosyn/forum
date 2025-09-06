document.addEventListener("DOMContentLoaded", () => {
    document.querySelectorAll(".category-bar a").forEach(link => {
        link.addEventListener("click", function (e) {
            e.preventDefault();
            const category = this.getAttribute("href").split("=")[1];

            fetch(category ? `/?category=${category}` : "/")
                .then(response => response.text())
                .then(html => {
                    document.querySelector("#post-list").innerHTML = new DOMParser()
                        .parseFromString(html, "text/html")
                        .querySelector("#post-list").innerHTML;
                })
                .catch(error => console.error("Error:", error));
        });
    });
});
