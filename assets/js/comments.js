document.addEventListener("DOMContentLoaded", () => {
    const commentForm = document.getElementById("comment-form");

    if (commentForm) {
        commentForm.addEventListener("submit", function (e) {
            e.preventDefault();
            console.log("Comment form submitted.");

            const formData = new FormData(this);

            fetch("/submitComment", {
                method: "POST",
                body: formData
            })
            .then(response => {
                if (response.ok) {
                    window.location.reload();
                } else {
                    throw new Error("Comment submission failed");
                }
            })
            .catch(error => {
                console.error("Error:", error);
                alert("Failed to submit comment. Please try again.");
            });
        });
    } else {
        console.log("No comment form found.");
    }
});
