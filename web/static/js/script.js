// Toggle sidebar on desktop
document.getElementById('toggle-sidebar-btn').addEventListener('click', function() {
    const sidebar = document.getElementById('sidebar');
    const contentArea = document.getElementById('main-content');
    const schoolName = document.getElementById('school-name');
    const navTexts = document.querySelectorAll('.nav-text');
    const icon = this.querySelector('i');

    if (sidebar.classList.contains('sidebar-collapsed')) {
        // Expand sidebar
        sidebar.classList.remove('sidebar-collapsed');
        contentArea.classList.remove('collapsed');
        schoolName.classList.remove('hidden');
        navTexts.forEach(text => text.classList.remove('hidden'));
        icon.classList.remove('fa-chevron-right');
        icon.classList.add('fa-chevron-left');
    } else {
        // Collapse sidebar
        sidebar.classList.add('sidebar-collapsed');
        contentArea.classList.add('collapsed');
        schoolName.classList.add('hidden');
        navTexts.forEach(text => text.classList.add('hidden'));
        icon.classList.remove('fa-chevron-left');
        icon.classList.add('fa-chevron-right');
    }
});

// Mobile menu toggle
document.getElementById('mobile-menu-btn').addEventListener('click', function() {
    const sidebar = document.getElementById('sidebar');
    const contentArea = document.getElementById('main-content');

    if (sidebar.classList.contains('mobile-open')) {
        sidebar.classList.remove('mobile-open');
        contentArea.classList.remove('mobile-shifted');
    } else {
        sidebar.classList.add('mobile-open');
        contentArea.classList.add('mobile-shifted');
    }
});

// Set active navigation item
function setActiveNav(clickedElement) {
    document.querySelectorAll('.nav-item').forEach(item => {
        item.classList.remove('active');
    });
    clickedElement.classList.add('active');

    // Update page title
    const navText = clickedElement.querySelector('.nav-text');
    if (navText) {
        document.getElementById('page-title').textContent = navText.textContent;
    }
}

// Close mobile sidebar when clicking outside
document.addEventListener('click', function(event) {
    const sidebar = document.getElementById('sidebar');
    const mobileMenuBtn = document.getElementById('mobile-menu-btn');

    if (window.innerWidth <= 768 && sidebar.classList.contains('mobile-open')) {
        const isClickInsideSidebar = sidebar.contains(event.target);
        const isClickOnMenuBtn = mobileMenuBtn.contains(event.target);

        if (!isClickInsideSidebar && !isClickOnMenuBtn) {
            sidebar.classList.remove('mobile-open');
            document.getElementById('main-content').classList.remove('mobile-shifted');
        }
    }
});

// Add fade-in effect to new content
document.body.addEventListener('htmx:afterOnLoad', function(evt) {
    if (evt.detail.target.id === 'main-content') {
        evt.detail.target.classList.add('fade-in');
        setTimeout(() => {
            evt.detail.target.classList.remove('fade-in');
        }, 300);
    }
});