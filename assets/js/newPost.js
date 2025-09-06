
    document.addEventListener('DOMContentLoaded', () => {
        const checkboxes = document.querySelectorAll('input[name="category"]');
        const form = document.getElementById('new-post-form');
        
        checkboxes.forEach(checkbox => {
            checkbox.addEventListener('change', () => {
                const selectedCount = document.querySelectorAll('input[name="category"]:checked').length;
                
                // Disable unchecked checkboxes when 3 are selected
                checkboxes.forEach(box => {
                    if (!box.checked) {
                        box.disabled = selectedCount >= 3;
                    }
                });
            });
        });

        // Form submission handler
        form.addEventListener('submit', (e) => {
            const selectedCategories = document.querySelectorAll('input[name="category"]:checked');

            if (selectedCategories.length === 0) {
                e.preventDefault();
                alert('Please select at least one category.');
                return;
        }

            const postData = {
                title: document.getElementById('title').value,
                categories: selectedCategories,
                content: document.getElementById('content').value
            };

            if (title.length < 5 || content.length < 10) {
                e.preventDefault();
                alert('Title must be at least 5 characters and content at least 10 characters long.');
                return;
            }

            console.log('Post Data:', postData);
            alert('Post created successfully!');
            // Handle form submission logic (e.g., send data to the server)
        });
    });
