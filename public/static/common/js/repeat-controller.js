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

    window.Ponzu.RepeatController = function(fieldName, parentSelector, childSelector, numberOfItems = 0) {
        let size = numberOfItems, n = numberOfItems;
        const itemLengthInputName = `${parentSelector}.length`;
        const removedItemsInputName = `${parentSelector}.removed`;

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
}());