// Basic client-side validation and functionality
document.addEventListener('DOMContentLoaded', function () {
    // Example: Add form validation if needed
    const forms = document.querySelectorAll('form');
    forms.forEach(form => {
        form.addEventListener('submit', function (event) {
            const requiredFields = form.querySelectorAll('[required]');
            let valid = true;

            requiredFields.forEach(field => {
                if (!field.value.trim()) {
                    valid = false;
                    field.classList.add('error');
                } else {
                    field.classList.remove('error');
                }
            });

            if (!valid) {
                event.preventDefault();
                alert('Please fill in all required fields.');
            }
        });
    });
});