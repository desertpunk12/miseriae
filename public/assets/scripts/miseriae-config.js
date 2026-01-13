
(function() {
    tailwind.config = {
        darkMode: "class",
        theme: {
            extend: {
                colors: {
                    "primary": "#ee2b8c",
                    "primary-light": "#ff8ec6",
                    "background-light": "#fff5f9", // Very light pinkish white
                    "background-dark": "#221019",
                    "text-main": "#181114",
                    "text-muted": "#896175",
                },
                fontFamily: {
                    "display": ["Spline Sans", "sans-serif"]
                },
                borderRadius: {"DEFAULT": "1rem", "lg": "1.5rem", "xl": "2.5rem", "2xl": "3.5rem", "full": "9999px"},
                boxShadow: {
                    'kawaii': '4px 4px 0px 0px #ee2b8c',
                    'kawaii-hover': '6px 6px 0px 0px #ee2b8c',
                }
            },
        },
    }
})();
