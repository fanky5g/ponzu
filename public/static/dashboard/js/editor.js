window.addEventListener("load", () => {
    const contentActiveClassName = "editor-content--active";
    const tabBar = document.querySelector(
        '.form-view-root .mdc-tab-bar[role="tablist"]',
    );
    const contentViews = document.querySelectorAll(
        ".form-view-content .editor-content",
    );

    if (tabBar) {
        const mdcTabBar = new mdc.tabBar.MDCTabBar(tabBar);

        window.addEventListener("MDCTabBar:activated", function(event) {
            event.preventDefault();

            const activeIndex = event.detail.index;
            for (let i = 0; i < contentViews.length; i++) {
                contentViews[i].classList.remove(contentActiveClassName);
            }

            const activeContentView = contentViews[activeIndex];
            if (activeContentView) {
                activeContentView.classList.add(contentActiveClassName);
            }
        });

        mdcTabBar.focusOnActivate = false;
        mdcTabBar.activateTab(0);
    } else {
        for (let i = 0; i < contentViews.length; i++) {
            contentViews[i].classList.add(contentActiveClassName);
        }
    }

    // initialize text fields
    const inputElements = document.querySelectorAll(
        ".form-view-root .mdc-text-field",
    );
    for (let i = 0; i < inputElements.length; i++) {
        new mdc.textField.MDCTextField(inputElements[i]);
    }

    // initialize checkboxes
    const checkboxElements = document.querySelectorAll(
        ".form-view-root .mdc-form-field",
    );
    for (let i = 0; i < checkboxElements.length; i++) {
        new mdc.checkbox.MDCCheckbox(checkboxElements[i]);
    }
});
