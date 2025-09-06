document.addEventListener("DOMContentLoaded", function() {
    // Restore last active tab
    const activeTab = localStorage.getItem("activeProfileTab");
    if (activeTab) {
        let tabElement = document.querySelector(`[data-bs-target="${activeTab}"]`);
        if (tabElement) {
            tabElement.click(); // Switch to the saved tab
        }
    }

    // Save active tab on click
    document.querySelectorAll("#profileTabs button").forEach(tab => {
        tab.addEventListener("click", function() {
            localStorage.setItem("activeProfileTab", this.dataset.bsTarget);
        });
    });
});