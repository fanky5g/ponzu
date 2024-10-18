const rowSelectionInputs = ['form[name="delete-item-form"] input[name="ids"]'];

const updateRowSelectionInputs = (value) => {
    for (let i = 0; i < rowSelectionInputs.length; i++) {
        const inputElement = document.querySelector(rowSelectionInputs[i]);
        if (inputElement) {
            inputElement.value = value;
        }
    }
};

const toggleCSVExportButtonEnabled = (hasSelectedItems) => {
    csvExportButton = document.querySelector('button[name="csv-export"]');
    if (csvExportButton) {
        csvExportButton.disabled = hasSelectedItems;
    }
};

const toggleDeleteButtonEnabled = (hasSelectedItems) => {
    const deleteItemsButton = document.querySelector(
        'button[name="delete-items"]',
    );
    if (deleteItemsButton) {
        deleteItemsButton.disabled = !hasSelectedItems;
    }
};

const triggerActions = (selectedItems) => {
    const hasSelectedItems = selectedItems.length > 0;

    updateRowSelectionInputs(selectedItems.join(","));
    toggleCSVExportButtonEnabled(hasSelectedItems);
    toggleDeleteButtonEnabled(hasSelectedItems);
};

const deleteItems = () => {
    return confirm(
        "[Ponzu] Please confirm:\n\nAre you sure you want to delete selected items?\nThis cannot be undone",
    );
};

const onSearchInputFocus = () => {
    const searchBarInnerRef = document.querySelector(
        ".search-bar .search-bar-inner",
    );
    if (searchBarInnerRef) {
        searchBarInnerRef.classList.add("open");
    }
};

const onSearchInputLostFocus = () => {
    const searchBarInnerRef = document.querySelector(
        ".search-bar .search-bar-inner",
    );
    if (searchBarInnerRef) {
        searchBarInnerRef.classList.remove("open");
    }
};

const rowsPerPageSelect = document.querySelector(
    ".mdc-data-table__pagination .mdc-select",
);

(function() {
    const items = document.querySelectorAll(".table tr.table-row");
    const hasItems = items?.length > 0;
    const tableView = document.querySelector(".mdc-data-table");

    if (tableView && hasItems) {
        const dataTable = new mdc.dataTable.MDCDataTable(tableView);
        window.addEventListener("MDCDataTable:selectedAll", function() {
            triggerActions(dataTable.getSelectedRowIds());
        });

        window.addEventListener("MDCDataTable:unselectedAll", function() {
            triggerActions([]);
        });

        window.addEventListener("MDCDataTable:rowSelectionChanged", function() {
            triggerActions(dataTable.getSelectedRowIds());
        });
    }

    const rowsPerPageList = document.querySelector(
        ".mdc-data-table__pagination .mdc-list",
    );

    if (rowsPerPageSelect) {
        new mdc.select.MDCSelect(rowsPerPageSelect);
    }

    if (rowsPerPageList) {
        const mdcList = new mdc.list.MDCList(rowsPerPageList);
        const elements = mdcList.listElements;
        const searchParams = new URLSearchParams(location.search);
        rowsPerPageList.addEventListener("MDCList:action", function(event) {
            const selected = elements[event.detail.index];
            const rowsPerPage = selected.dataset.value;
            searchParams.set("count", rowsPerPage);
            searchParams.set("offset", "0");
            window.location.replace(
                `${location.pathname}?${searchParams.toString()}`,
            );
        });
    }
})();
