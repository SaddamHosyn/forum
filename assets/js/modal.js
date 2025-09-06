// js/modal.js
document.addEventListener('DOMContentLoaded', function () {
    // Open modals when clicking nav links
    document.addEventListener('click', function (e) {
        if (e.target.id === 'nav-login') {
            e.preventDefault();
            document.getElementById('login-modal').style.display = 'flex';
        } else if (e.target.id === 'nav-register') {
            e.preventDefault();
            document.getElementById('register-modal').style.display = 'flex';
        }
    });

    // Close modals
    document.querySelectorAll('.close-modal').forEach(button => {
        button.addEventListener('click', () => {
            const modalId = button.getAttribute('data-modal');
            document.getElementById(modalId).style.display = 'none';
        });
    });

    // Switch between modals for all elements with the class
    document.querySelectorAll('.switch-to-register').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            document.getElementById('login-modal').style.display = 'none';
            document.getElementById('register-modal').style.display = 'flex';
        });
    });

    document.querySelectorAll('.switch-to-login').forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            document.getElementById('register-modal').style.display = 'none';
            document.getElementById('login-modal').style.display = 'flex';
        });
    });

    // Close modals when clicking outside
    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            e.target.style.display = 'none';
        }
    });
});

