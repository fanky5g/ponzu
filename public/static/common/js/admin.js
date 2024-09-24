const actionsMenu = document.querySelector(".app-bar-right .mdc-menu");
if (actionsMenu) {
    const menu = new mdc.menu.MDCMenu(actionsMenu);
    const appBarActionsButton = document.querySelector(".app-bar-right .actions");

    appBarActionsButton.addEventListener("click", function() {
        menu.open = true;
    });
}
