window.Ponzu.initializeNestedRepeater = (
    entityName,
    selector,
    cloneSelector,
    positionalPlaceholder,
    template,
    numItems,
) => {
    let n = numItems;
    let size = numItems;
    const parser = new DOMParser();

    const controlSelector = `.__ponzu-repeat.${window.Ponzu.cleanQueryPath(selector)}`;
    const nestedRepeaterControl = document.querySelector(controlSelector);
    const getChildren = function() {
        return nestedRepeaterControl.querySelectorAll(cloneSelector);
    };

    const setHiddenInputValue = (name, value) => {
        const selectorName= `.__ponzu-repeat.${selector}.${name}`;
        let input = nestedRepeaterControl.querySelector(`input[name="${selectorName}"]`);
        if (!input) {
            input = parser.parseFromString(`<input name="${selectorName}" type="hidden" value="${value}" />`, "text/html").body.firstChild;
            nestedRepeaterControl.appendChild(input);
            return;
        }

        input.value = value;
    };

    const insertEmptyListHelperText = () => {
        const emptyListHelperText = parser.parseFromString(
            `<label class="empty-text">No ${entityName} added.</label>`,
            "text/html",
        ).body.firstChild;

        const parentControls = nestedRepeaterControl.querySelector('.parent.nested-repeater-controls');
        if (parentControls) {
            $(emptyListHelperText).insertBefore($(parentControls));
        } else {
            $(nestedRepeaterControl).append(emptyListHelperText);
        }
    };

    const removeEmptyListHelperText = () => {
        nestedRepeaterControl.querySelector('label.empty-text')?.remove();
    };

    const createChild = function(i) {
        removeEmptyListHelperText();

        const childTemplate = template.replace(new RegExp(positionalPlaceholder, 'g'), i);

        const childNode = $('<div/>').html(childTemplate).contents();
        const parentControls = nestedRepeaterControl.querySelector('.parent.nested-repeater-controls');
        if (parentControls) {
            childNode.insertBefore($(parentControls));
        } else {
            $(nestedRepeaterControl).append(childNode);
        }

        const controls = createChildControls(i);
        $(childNode).append(controls);
        n++; size++;
        setHiddenInputValue("length", size);

        window.Ponzu.initializeEditorControls(controlSelector);
    };

    const deleteChild = function (el) {
        const wrapper = $(el).parent().closest(cloneSelector);
        wrapper.remove();

        size = size - 1;
        setHiddenInputValue("length", size);
        if (!size) {
            insertEmptyListHelperText();
        }
    };

    const addRepeater = function (e) {
        e.preventDefault();
        createChild(n);
    };

    const delRepeater = function (e, elementIndex) {
        e.preventDefault();
        const selectorName= `.__ponzu-repeat.${selector}.removed`;
        let removedElementsInput = nestedRepeaterControl.querySelector(`input[name="${selectorName}"]`);

        if (!removedElementsInput) {
            removedElementsInput = parser.parseFromString(
                `<input name="${selectorName}" type="hidden" value="" />`,
                "text/html",
            ).body.firstChild;
            nestedRepeaterControl.appendChild(removedElementsInput);
        }

        deleteChild(e.target);

        const currentValue = removedElementsInput.value;
        removedElementsInput.value = [currentValue, elementIndex].filter(el => el.toString().trim() !== "").join(',');
    };

    const createChildControls = function (elementIndex) {
        const deleteButton = parser.parseFromString(`
                <button class="mdc-icon-button material-symbols-rounded">
                  <div class="mdc-icon-button__ripple"></div>
                  remove
                </button>
            `, "text/html").body.firstChild;

        deleteButton.addEventListener('click', function (e) {
            delRepeater(e, elementIndex);
        });

        const deleteButtonContainer = parser.parseFromString(
            `<div class="repeater-del"></div>`,
            "text/html"
        ).body.firstChild;

        deleteButtonContainer.appendChild(deleteButton);
        return deleteButtonContainer;
    };

    const createParentControls = function() {
        const container = parser.parseFromString(
            '<div class="parent nested-repeater-controls"></div>',
            "text/html",
        ).body.firstChild;

        const addButton = parser.parseFromString(`
                <button class="mdc-icon-button material-symbols-rounded">
                  <div class="mdc-icon-button__ripple"></div>
                  add_row_above
                </button>
            `, "text/html").body.firstChild;

        addButton.classList.add('add-block');
        addButton.addEventListener('click', addRepeater);

        container.appendChild(addButton);
        nestedRepeaterControl.appendChild(container);
    };

    const children = getChildren();
    children.forEach((childNode, index) => {
        childNode.appendChild(createChildControls(index));
    });

    if (!n) {
        insertEmptyListHelperText();
    }

    createParentControls();
};