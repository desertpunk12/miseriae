(function () {
    const toggleBtn = document.getElementById('theme-toggle');

    // Function to update local state if needed, but we rely mostly on CSS
    // mirroring the 'dark' class on document.documentElement

    if (toggleBtn) {
        toggleBtn.addEventListener('click', () => {
            const html = document.documentElement;
            const isDark = html.classList.contains('dark');

            if (isDark) {
                html.classList.remove('dark');
                localStorage.setItem('theme', 'light');
            } else {
                html.classList.add('dark');
                localStorage.setItem('theme', 'dark');
            }
        });
    }
})();
