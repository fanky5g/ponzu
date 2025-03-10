(function () {
  function eachNode(rootNode, callback) {
    if (!callback) {
      return;
    }

    if (callback(rootNode) === false) {
      return false;
    }

    if (rootNode.hasChildNodes()) {
      for (const node of rootNode.childNodes) {
        if (eachNode(node, callback) === false) {
          return;
        }
      }
    }
  }

  const get = function (object, path) {
    path = path.replace(/\[(\w+)]/g, ".$1");
    path = path.replace(/^\./, "");

    let keys = path.split(".");
    for (let i = 0, n = keys.length; i < n; ++i) {
      let key = keys[i];

      if (key in object) {
        object = object[key];
      } else {
        return;
      }
    }

    return object;
  };

  const newJsonValueReplacer = (object) => (_, path) => {
    return get(object, path);
  };

  function processTextNode(node, pattern, replacer) {
    if (pattern instanceof RegExp && node.textContent.match(pattern)) {
      node.textContent = node.textContent.replaceAll(pattern, replacer);
    }
  }

  function processAttr(attr, pattern, replacer) {
    if (pattern instanceof RegExp && attr.value.match(pattern)) {
      attr.value = attr.value.replace(pattern, replacer);
    }
  }

  const publicPath = window.Ponzu?.publicPath ?? "/";
  const request = (method, url) => {
    const { promise, resolve } = Promise.withResolvers();
    const request = new XMLHttpRequest();

    request.addEventListener("load", () => {
      resolve(request.response.data);
    });

    request.open(method, url);
    request.responseType = "json";
    request.send();
    return promise;
  };

  const loadOptions = async (
    contentType,
    count = 20,
    offset = 0,
    order = "DESC",
  ) =>
    request(
      "GET",
      `${publicPath}/api/references?type=${contentType}&count=${count}&offset=${offset}&order=${order}`,
    );

  const loadOption = async (contentType, id) =>
    request("GET", `${publicPath}/api/references/${id}?type=${contentType}`);

  const appendRow = (node, row) => {
    node.appendChild(row);
  };

  const prependRow = (node, row) => {
    if (!node.firstChild) {
      return appendRow(node, row);
    }

    node.insertBefore(row, node.firstChild);
  };

  const newNodeReplacerFromJSONValue = (node, value) => {
      const replacer = newJsonValueReplacer(value);

      if (node.nodeType === Node.TEXT_NODE) {
        processTextNode(node, /"@>(.*)"/gm, replacer);
        return;
      }

      if (node.nodeType === Node.ELEMENT_NODE && node.hasAttributes()) {
        for (const attr of node.attributes) {
          processAttr(attr, /@>(.*)/, replacer);
        }
      }
  };

  // renderOptionsIntoNode renders passed options into HTML Node.
  // @param mode = enum(PREPEND | APPEND)
  const parser = new DOMParser();
  async function renderOptionsIntoNode(
    options,
    template,
    node,
    mode = "APPEND",
    onNodeInsertCallback = () => {},
  ) {
    const rows = options.map((option) => {
      const optionRow = parser.parseFromString(template, "text/html");
      eachNode(optionRow, (node) => newNodeReplacerFromJSONValue(node, option));
      return optionRow;
    });

    const insertRow = mode === "APPEND" ? appendRow : prependRow;
    for (let row of rows) {
      const childNode = row.body.firstChild;
      insertRow(node, childNode);
      onNodeInsertCallback(childNode);
    }
  }

  const createOptionsLoader = (contentType, pageSize = 20, order = "DESC") => {
    let offset = 0;
    let hasMore = true;

    return async (filterFunc = () => true) => {
      if (!hasMore) {
        return [];
      }

      const { references, size } = await loadOptions(
        contentType,
        pageSize,
        offset,
        order,
      );
      offset = offset + references.length;
      if (offset === size) {
        hasMore = false;
      }

      return references.filter(filterFunc);
    };
  };

  const createSingleSelect = (contentType, selector, rowTemplate, loadOptions = async () => {}) => {
    const select = document.querySelector(`.mdc-select.single-select.${window.Ponzu.cleanQueryPath(selector)}`);
    const initialValue = select.querySelector('input[type="hidden"]')?.value;
    let mdcSelect = new mdc.select.MDCSelect(select);
    const mdcMenu = mdcSelect.menu;

    let filterFunc = () => true;
    if (initialValue) {
      filterFunc = (option) => {
        return option.id !== initialValue;
      };
    }

    async function update() {
      const options = await loadOptions(filterFunc);
      if (options.length) {
        await renderOptionsIntoNode(options, rowTemplate, mdcMenu.list.root);
        mdcSelect = new mdc.select.MDCSelect(select);
      }
    }

    return {
      initialize: async function () {
        const hasValidOptions = mdcMenu.items.some((item) => {
          return Boolean(item.dataset.value);
        });

        if (!hasValidOptions) {
          await update();
          if (initialValue && !mdcSelect.value) {
            const option = await loadOption(contentType, initialValue);
            await renderOptionsIntoNode(
              [option],
              rowTemplate,
              mdcMenu.list.root,
              "PREPEND",
            );

            mdcSelect = new mdc.select.MDCSelect(select);
          }
        }
      },
      registerScrollHandlers: function () {
        mdcMenu.root.addEventListener("scroll", async function (e) {
          const element = e.target;
          const scrollHeight = element.scrollHeight - element.clientHeight;
          const hasScrolledToBottom = element.scrollTop === scrollHeight;

          if (hasScrolledToBottom) {
            await update();
          }
        });
      },
    };
  };

  const createMultiSelect = async (contentType, selector, rowTemplate, selectedItemTemplate, loadOptions = async () => {}) => {
    const parentSelector = `.__ponzu-repeat.${window.Ponzu.cleanQueryPath(selector)}`;
    const childSelector = `.mdc-chip`;

    const select = document.querySelector(`.mdc-select.multi-select.${window.Ponzu.cleanQueryPath(selector)}`);
    const chipContainer = select.querySelector('div.mdc-chip-set');
    const selectedOptions = [];

    let mdcSelect = new mdc.select.MDCSelect(select);
    const mdcMenu = mdcSelect.menu;

    const repeatController = window.Ponzu.RepeatController(
        selector,
        parentSelector,
        childSelector,
        chipContainer.querySelectorAll(childSelector).length,
    );

    const onRemoveChip = (event) => {
      event.stopPropagation();

      const chip = event.target.offsetParent;
      const value = chip.dataset.value;

      if (value === mdcSelect.value) {
        mdcSelect.setSelectedIndex(-1);
      }

      selectedOptions.splice(selectedOptions.indexOf(value), 1);
      repeatController.onChildRemoved(chip);

      chip.parentNode.removeChild(chip);
    };

    const onChildAdd = (childNode) => {
      const removeButton = childNode.querySelector('i[role="button"]');
      removeButton?.addEventListener('click', onRemoveChip);

      repeatController.onChildAdded(childNode);
    };

    const initializeSelect = (reset = false) => {
      if (reset && mdcSelect) {
          mdcSelect.destroy();
          mdcSelect = new mdc.select.MDCSelect(select);
      }

      mdcSelect.listen('MDCSelect:change', async (event) => {
        const { detail: { value } } = event;

        if (value && !selectedOptions.includes(value)) {
          selectedOptions.push(value);
          const option = await loadOption(contentType, value);
          await renderOptionsIntoNode(
              [option],
              selectedItemTemplate,
              chipContainer,
              "APPEND",
              onChildAdd
          );
        }
      });
    };

    let filterFunc = () => true;
    const initialValue = "";
    if (initialValue) {
        filterFunc = (option) => {
          return  selectedOptions.indexOf(option.id) === -1;
        };
    }

    async function update() {
      const options = await loadOptions(filterFunc);
      if (options.length) {
        await renderOptionsIntoNode(options, rowTemplate, mdcMenu.list.root, onChildAdd);
        initializeSelect(true);
      }
    }

    return {
      initialize: async function () {
        initializeSelect();

        const chipNodes = chipContainer.querySelectorAll(childSelector);
        const chipNodePromises = [];

        chipNodes.forEach((chipNode) => {
          const removeButton = chipNode.querySelector('i[role="button"]');
          removeButton?.addEventListener('click', onRemoveChip);

          chipNodePromises.push(new Promise((resolve, reject) => {
            const value = chipNode.dataset.value;
            selectedOptions.push(value);

            loadOption(contentType, value).then(referencedData => {
              eachNode(chipNode, (node) => newNodeReplacerFromJSONValue(node, referencedData));
              resolve();
            }).catch(reject);
          }));
        });

        await Promise.all(chipNodePromises);

        const hasValidOptions = mdcMenu.items.some((item) => {
          return Boolean(item.dataset.value);
        });

        if (!hasValidOptions) {
          await update();
        }
      },
      registerScrollHandlers: function () {
        mdcMenu.root.addEventListener("scroll", async function (e) {
          const element = e.target;
          const scrollHeight = element.scrollHeight - element.clientHeight;
          const hasScrolledToBottom = element.scrollTop === scrollHeight;

          if (hasScrolledToBottom) {
            await update();
          }
        });
      },
    };
  };

  window.Ponzu.initializeReferenceLoader = async (
      contentType,
      selector,
      selectType,
      optionTemplate,
      selectedOptionTemplate,
  ) => {
    let select;

    switch (selectType) {
      case "single":
        select = createSingleSelect(
            contentType,
            selector,
            optionTemplate,
            createOptionsLoader(contentType),
        );
        break;
      case "multiple":
        select = await createMultiSelect(
              contentType,
              selector,
              optionTemplate,
              selectedOptionTemplate,
              createOptionsLoader(contentType),
        );
        break;
      default:
        console.error(`Invalid select type: ${selectType}`);
        break;
    }

    await select.initialize();
    select.registerScrollHandlers();
  };
})();
