(function () {
    const ErrValueRequired = new Error("[Repeat Controller] Expected child node value to be defined");
    const parser = new DOMParser();

    const setHiddenInputValue = function(parentNode, name, value, mode = "REPLACE") {
        let input = parentNode.querySelector(`input[name="${name}"]`);

        if (!input) {
            const inputTemplate = `<input name="${name}" type="hidden" value="" />`;
            const inputNode = parser.parseFromString(inputTemplate, "text/html");
            input = inputNode.body.firstChild;
            parentNode.appendChild(inputNode.body.firstChild);
        }

        if (mode === "REPLACE") {
            input.value = value;
            return;
        }

        input.value = input.value ?  [input.value, value].join(", "): value;
    };

    window.Ponzu.RepeatController = function(fieldName, parentSelector, numberOfItems = 0) {
        let size = numberOfItems, n = numberOfItems;
        const normalizedParentSelector = window.Ponzu.normalizeQueryPath(parentSelector);
        const itemLengthInputName = `${normalizedParentSelector}.length`;
        const removedItemsInputName = `${normalizedParentSelector}.removed`;

        const parent = document.querySelector(parentSelector);
        setHiddenInputValue(parent, itemLengthInputName, numberOfItems);

        return {
            onChildAdded: (childNode) => {
                const value = childNode.dataset.value;
                if (!value) {
                    console.error(ErrValueRequired);
                    return;
                }

                const insertedChildInputName = `${fieldName}.${n}`;

                size = size + 1;
                n = n + 1;

                setHiddenInputValue(childNode, insertedChildInputName, value);
                setHiddenInputValue(parent, itemLengthInputName, size);
            },
            onChildRemoved: (childNode) => {
                size = size - 1;

                const value = childNode.dataset.value;
                if (!value) {
                    console.error(ErrValueRequired);
                    return;
                }

                const inputNode = childNode.querySelector(`input[value="${value}"]`);
                const removedItemIndex = inputNode.name.replace(`${fieldName}.`, "");

                setHiddenInputValue(parent, itemLengthInputName, size);
                setHiddenInputValue(parent, removedItemsInputName, removedItemIndex, "APPEND");
            },
        };
    };

    window.Ponzu.initializeRepeatControl = function(
        entityName,
        selector,
        cloneSelector,
        positionalPlaceholder,
        template,
        numItems = 0,
        elementType = 'input') {
        let n = numItems;
        let size = numItems;
        const parser = new DOMParser();

        const controlSelector = `.__ponzu-repeat.${window.Ponzu.cleanQueryPath(selector)}`;
        const repeaterControl = document.querySelector(controlSelector);

        const getChildren = function() {
            return repeaterControl.querySelectorAll(cloneSelector);
        };

        const removeEmptyListHelperText = () => {
            repeaterControl.querySelector(':scope > label.empty-text')?.remove();
        };

        const setHiddenInputValue = (name, value) => {
            const selectorName= `.__ponzu-repeat.${selector}.${name}`;
            let input = repeaterControl.querySelector(`:scope > input[name="${selectorName}"]`);
            if (!input) {
                input = parser.parseFromString(`<input name="${selectorName}" type="hidden" value="${value}" />`, "text/html").body.firstChild;
                repeaterControl.appendChild(input);
                return;
            }

            input.value = value;
        };

        const insertEmptyListHelperText = () => {
            const emptyListHelperText = parser.parseFromString(
                `<label class="empty-text">No ${entityName} added.</label>`,
                "text/html",
            ).body.firstChild;

            const parentControls = repeaterControl.querySelector(`:scope > .parent.${elementType}-repeater-controls`);
            if (parentControls) {
                $(emptyListHelperText).insertBefore($(parentControls));
            } else {
                $(repeaterControl).append(emptyListHelperText);
            }
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

        const delRepeater = function (e, elementIndex) {
            e.preventDefault();
            const selectorName= `.__ponzu-repeat.${selector}.removed`;
            let removedElementsInput = repeaterControl.querySelector(`:scope > input[name="${selectorName}"]`);

            if (!removedElementsInput) {
                removedElementsInput = parser.parseFromString(
                    `<input name="${selectorName}" type="hidden" value="" />`,
                    "text/html",
                ).body.firstChild;
                repeaterControl.appendChild(removedElementsInput);
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

        const createChild = function(i) {
            removeEmptyListHelperText();

            const childTemplate = template.replace(new RegExp(positionalPlaceholder, 'g'), i);

            const childNode = $('<div/>').html(childTemplate).contents();
            const parentControls = repeaterControl.querySelector(`:scope > .parent.${elementType}-repeater-controls`);
            if (parentControls) {
                childNode.insertBefore($(parentControls));
            } else {
                $(repeaterControl).append(childNode);
            }

            const controls = createChildControls(i);
            $(childNode).append(controls);
            n++; size++;
            setHiddenInputValue("length", size);

            window.Ponzu.initializeEditorControls(controlSelector);
        };

        const addRepeater = function (e) {
            e.preventDefault();
            createChild(n);
        };

        const createParentControls = function() {
            const container = parser.parseFromString(
                `<div class="parent ${elementType}-repeater-controls"></div>`,
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
            repeaterControl.appendChild(container);
        };

        const children = getChildren();
        children.forEach((childNode, index) => {
            childNode.appendChild(createChildControls(index));
        });

        if (!n) {
            insertEmptyListHelperText();
        }

        createParentControls();
    }
}());