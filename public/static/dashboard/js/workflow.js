window.addEventListener("load", () => {
    const workflowActions = document.querySelector(".workflow-actions");
    const mdcMenu = workflowActions?.querySelector(".mdc-menu");

    if (mdcMenu) {
        const mdcMenuInstance = new mdc.menu.MDCMenu(mdcMenu);
        const button = workflowActions.querySelector("button");

        button?.addEventListener("click", function() {
            mdcMenuInstance.open = true;
        });
    }
});
