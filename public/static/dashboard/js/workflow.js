window.addEventListener("load", () => {
    const workflowActions = document.querySelector(".workflow-actions");
    const mdcMenu = workflowActions?.querySelector(".mdc-menu");
    const workflowState = workflowActions?.querySelector(
        'input[name="workflow_state"]',
    );

    if (mdcMenu) {
        const mdcMenuInstance = new mdc.menu.MDCMenu(mdcMenu);
        const workflowActionsList = mdcMenu.querySelector(".mdc-list");
        const mdcList = new mdc.list.MDCList(workflowActionsList);
        const button = workflowActions.querySelector("button");

        button?.addEventListener("click", function() {
            mdcMenuInstance.open = true;
        });

        workflowActionsList?.addEventListener("MDCList:action", function(event) {
            event.preventDefault();
            const selected = mdcList.listElements[event.detail.index];
            workflowState.value = selected.dataset.value;
            workflowState.form.action = `${window.Ponzu.publicPath}/edit/workflow`;
            workflowState.form.requestSubmit();
        });
    }
});
