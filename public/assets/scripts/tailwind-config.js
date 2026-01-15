// Universal Tailwind Configuration
// Merges behaviors from miseriae-config.js, blog-tailwind-config.js, and tailwind-config.js

(function () {
  // 1. Initialize Theme (Dark/Light Mode)
  // Check localStorage or system preference to prevent FOUC (Flash of Unstyled Content)
  try {
    var storedTheme = localStorage.getItem('theme');
    var prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
    if (storedTheme === 'dark' || (storedTheme === null && prefersDark)) {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  } catch (e) {
    console.error('Theme initialization failed:', e);
  }

  // 2. Detect Context (Blog vs Main Site)
  // The Blog and Cosplays sections use the "Bubblegum Pop" theme
  const isBubblegum = window.location.pathname.startsWith('/blog') || window.location.pathname.startsWith('/cosplays');

  // 3. Define Palettes
  const mainColors = {
    primary: '#ee2b8c',              // Magenta
    'primary-light': '#ff8ec6',      // Light Magenta
    'primary-hover': '#d41b76',      // Darker Magenta for hovers
    'background-light': '#fff5f9',   // Very light pinkish white
    'background-dark': '#221019',    // Deep dark
    'text-main': '#181114',
    'text-muted': '#896175'
  };

  const bubblegumColors = {
    primary: '#ff69b4',              // Hot Pink
    'primary-light': '#ffd1e8',      // Light Pink
    'primary-dark': '#e05da5',       // Darker Hot Pink
    'background-light': '#fff5f9',   // Same as main
    'background-dark': '#3c1e30',    // Distinct dark mode
    'accent-pink': '#f786c2',
    'accent-purple': '#e0b0ff',
    'text-dark': '#4a1936'           // Specific text color
  };

  // 4. Merge Configuration
  // We use spread syntax to merge, with the active theme coming last to override collisions (like 'primary')
  const activeColors = isBubblegum ?
    { ...mainColors, ...bubblegumColors } :
    { ...bubblegumColors, ...mainColors };

  // 5. Define Shared and Specific Extended Config
  const config = {
    darkMode: 'class',
    theme: {
      extend: {
        colors: activeColors,
        fontFamily: {
          display: ['Spline Sans', 'sans-serif'],
          body: ['Noto Sans', 'sans-serif']
        },
        borderRadius: {
          DEFAULT: '1rem',
          lg: '1.5rem',
          // Context-aware radii
          xl: isBubblegum ? '2rem' : '2.5rem',
          '2xl': isBubblegum ? '3rem' : '3.5rem',
          '3xl': '4rem',
          full: '9999px',
          // Bubblegum specific shapes
          'bubble-sm': '30% 70% 70% 30% / 50% 50% 50% 50%',
          'bubble-md': '60% 40% 30% 70% / 50% 40% 60% 50%',
          'bubble-lg': '40% 60% 60% 40% / 70% 50% 50% 30%',
        },
        backgroundImage: {
          'sparkles': 'url("data:image/svg+xml,%3Csvg width=\'40\' height=\'40\' viewBox=\'0 0 40 40\' xmlns=\'http://www.w3.org/2000/svg\'%3E%3Cg fill=\'none\' fill-rule=\'evenodd\'%3E%3Cg fill=\'%23ff69b4\' fill-opacity=\'0.15\'%3E%3Cpath d=\'M20 0L22.5 7.5L30 10L22.5 12.5L20 20L17.5 12.5L10 10L17.5 7.5L20 0Z\'/%3E%3Cpath d=\'M5 15L7.5 22.5L15 25L7.5 27.5L5 35L2.5 27.5L-5 25L2.5 22.5L5 15Z\'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E")',
          'gradient-pop': 'linear-gradient(135deg, #ff69b4, #f786c2, #e0b0ff)',
        },
        boxShadow: {
          // Main Shadows
          'kawaii': '4px 4px 0px 0px #ee2b8c',
          'kawaii-hover': '6px 6px 0px 0px #ee2b8c',
          'soft': '0 4px 20px -2px rgba(238, 43, 140, 0.1)',
          'hover': '0 10px 25px -5px rgba(238, 43, 140, 0.15)',
          // Bubblegum Shadows
          'pop': '0 10px 20px -5px rgba(255, 105, 180, 0.3), 0 4px 6px -2px rgba(255, 105, 180, 0.1)',
          'pop-lg': '0 20px 30px -10px rgba(255, 105, 180, 0.4), 0 8px 12px -4px rgba(255, 105, 180, 0.15)',
        }
      }
    }
  };

  // 6. Apply to Tailwind
  if (typeof tailwind !== 'undefined') {
    tailwind.config = config;
  } else {
    window.tailwind = window.tailwind || {};
    window.tailwind.config = config;
  }
})();
