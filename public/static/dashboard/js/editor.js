(function() {
    const initializeSelectControls = (parentSelector) => {
        const selectControlSelector = [parentSelector, 'div.mdc-select'].filter(Boolean).join(' ');
        const selectControls = document.querySelectorAll(selectControlSelector);

        selectControls.forEach(selectControl => {
            new mdc.select.MDCSelect(selectControl);
        });
    };

    const initializeEditorControls = (parentSelector) => {
        const contentActiveClassName = "editor-content--active";
        const tabBar = document.querySelector(
            '.form-view-root .mdc-tab-bar[role="tablist"]',
        );

        const contentViews = document.querySelectorAll(
            ".form-view-content .editor-content",
        );

        const exceptionBanner = document.querySelector(".mdc-banner.exception");
        if (exceptionBanner) {
            const mdcExceptionBanner = new mdc.banner.MDCBanner(exceptionBanner);
            mdcExceptionBanner.open();
        }

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

        // initialize select controls
        initializeSelectControls(parentSelector);
    };

    if (typeof window.Ponzu !== "undefined") {
        window.Ponzu.initializeEditorControls = initializeEditorControls;
        initializeEditorControls();
        return;
    }

    window.addEventListener("load", () => {
        window.Ponzu.initializeEditorControls = initializeEditorControls;
        initializeEditorControls();
    });
})();


