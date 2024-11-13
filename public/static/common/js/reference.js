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
    path = path.replace(/\[(\w+)\]/g, ".$1");
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
      `${publicPath}/api/content/?type=${contentType}&count=${count}&offset=${offset}&order=${order}`,
    );

  const loadOption = async (contentType, id) =>
    request("GET", `${publicPath}/api/content/${id}?type=${contentType}`);

  const appendRow = (node, row) => {
    node.appendChild(row);
  };

  const prependRow = (node, row) => {
    if (!node.firstChild) {
      return appendRow(node, row);
    }

    node.insertBefore(row, node.firstChild);
  };

  // renderOptionsIntoNode renders passed options into HTML Node.
  // @param mode = enum(PREPEND | APPEND)
  const parser = new DOMParser();
  async function renderOptionsIntoNode(
    options,
    template,
    node,
    mode = "APPEND",
  ) {
    const rows = options.map((option) => {
      const optionRow = parser.parseFromString(template, "text/html");
      const replacer = newJsonValueReplacer(option);

      eachNode(optionRow, (node) => {
        if (node.nodeType === Node.TEXT_NODE) {
          processTextNode(node, /"@>(.*)"/gm, replacer);
          return;
        }

        if (node.nodeType === Node.ELEMENT_NODE && node.hasAttributes()) {
          for (const attr of node.attributes) {
            processAttr(attr, /@>(.*)/, replacer);
          }
        }
      });

      return optionRow;
    });

    const insertRow = mode === "APPEND" ? appendRow : prependRow;
    for (let row of rows) {
      insertRow(node, row.body.firstChild);
    }
  }

  const createOptionsLoader = (contentType, pageSize = 20, order = "DESC") => {
    let offset = 0;
    let hasMore = true;

    return async (filterFunc = () => true) => {
      if (!hasMore) {
        return [];
      }

      const options = await loadOptions(contentType, pageSize, offset, order);
      offset = offset + options.length;
      if (options.length < pageSize) {
        hasMore = false;
      }

      return options.filter(filterFunc);
    };
  };

  const createSelectElement = (
    contentType,
    rowTemplate,
    loadOptions = async () => {},
  ) => {
    const select = document.querySelector(`.mdc-select.${contentType}`);
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

  window.Ponzu.initializeReferenceLoader = async (contentType, rowTemplate) => {
    const select = createSelectElement(
      contentType,
      rowTemplate,
      createOptionsLoader(contentType),
    );

    await select.initialize();
    select.registerScrollHandlers();
  };
})();
