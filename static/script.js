const SyntaxKeyord = "keyword";
const SyntaxString = "string";
const SyntaxFunction = "function";
const SyntaxComment = "comment";
const SyntaxPlain = "plain";
const SyntaxNumber = "number";

// javascript treat unquoted key as string, oh my...
// don't define key-value object like this:
// const Highlight = {SyntaxKeyord: "purple"};
const Highlight = {};
Highlight[SyntaxKeyord] = "purple";
Highlight[SyntaxString] = "rgb(196,27,22)";
Highlight[SyntaxPlain] = "black";
Highlight[SyntaxComment] = "gray";
Highlight[SyntaxFunction] = "rgb(62,117,122)"
Highlight[SyntaxNumber] = "blue";

// Go keyword and type
const Keywords = [
    "var", "func", "if", "else", "for", "range", "continue", "return",
    "int", "string", "nil",
];
const Delimiters = [' ', '(', ')', '{', '}', '.', '\t', '[', ']'];

// return [[syntaxType, token]]
function parse(line) {
    if (line === "") {
        return [];
    }

    let syntax = [];
    let token = "";
    let delim = "";
    let inString = false;
    for (let i = 0; i < line.length; i++) {
        if (!inString && line.slice(i, i + 2) === "//") {
            syntax.push([SyntaxComment, line.slice(i)])
            break
        }
        if (line[i] === '"') {
            if (!inString) {
                inString = true;
            } else {
                inString = false;
                syntax.push([SyntaxString, token.concat(line[i])]);
                token = "";
                continue
            }
        }
        if (inString || !Delimiters.includes(line[i])) {
            token = token.concat(line[i]);
            continue;
        }
        delim = line[i];
        if (token !== "") {
            let type = SyntaxPlain;
            if (Keywords.includes(token)) {
                type = SyntaxKeyord
            } else if (delim === '(') {
                type = SyntaxFunction;
            } else if (!isNaN(Number(token))) {
                type = SyntaxNumber;
            }
            syntax.push([type, token]);
        }
        syntax.push([SyntaxPlain, delim])
        token = "";
    }
    if (token !== "") {
        let type = SyntaxPlain;
        if (Keywords.includes(token)) {
            type = SyntaxKeyord
        } else if (token.startsWith('"') && token.endsWith('"')) {
            type = SyntaxString;
        }
        syntax.push([type, token]);
    }
    return syntax;
}

function render(syntaxType, token) {
    let span = document.createElement("span");
    span.append(token);
    span.style.color = Highlight[syntaxType];
    return span;
}

function tokenColor(token) {
    let type = SyntaxPlain;
    if (Keywords.includes(token)) {
        type = SyntaxKeyord
    } else if (token.startsWith('"') && token.endsWith('"')) {
        type = SyntaxString;
    }
    return Highlight[type]
}

function initEditor(src) {
    let lines = src.trim().split("\n");
    const gutter = document.getElementById("gutter")
    gutter.innerHTML = '';
    for (let i = 1; i <= lines.length; i++) {
        let div = document.createElement("div");
        div.append(i)
        gutter.append(div)
    }

    const editarea = document.getElementById("editarea");
    editarea.innerHTML = '';
    for (let i = 0; i < lines.length; i++) {
        let line = lines[i];
        let renderLine = document.createElement("div");
        if (line === "") {
            renderLine.append("\n")
        } else {
            let syntax = parse(line);
            // console.log(i, syntax);
            for (let j = 0; j < syntax.length; j++) {
                let span = render(syntax[j][0], syntax[j][1]);
                renderLine.append(span);
            }
        }
        // refresh syntax highlight on changes
        // const observer = new MutationObserver((mutations) => {
        //     for (const mutation of mutations) {
        //         const selection = window.getSelection();
        //         const range = selection.getRangeAt(0);
        //         range.endContainer.parentElement.style.color = tokenColor(mutation.target.textContent);
        //     }
        // });
        // observer.observe(renderLine, { characterData: true, subtree: true });
        editarea.append(renderLine);
    }

    // window.getSelection().setPosition(editarea.firstChild, 0);
    // editarea.addEventListener('keydown', function (event) {
    //     if (event.key === 'Tab') {
    //         // disable focus shift
    //         event.preventDefault();

    //         const selection = window.getSelection();
    //         const range = selection.getRangeAt(0);
    //         // Insert the tab character
    //         range.insertNode(document.createTextNode("\t"));
    //         // Advance the caret position 
    //         range.setStart(range.endContainer, range.endOffset);
    //         range.collapse(false);
    //     }
    // });

    // editarea.addEventListener("paste", function (event) {
    //     event.preventDefault();
    //     const pastedText = event.clipboardData.getData("text/plain");
    //     const range = window.getSelection().getRangeAt(0);
    //     range.deleteContents();
    //     range.insertNode(document.createTextNode(pastedText));
    //     // TODO: paste multiple line
    //     range.collapse(false);
    // });

    // const mutationObserver = new MutationObserver(function (mutations) {
    //     for (let mutation of mutations) {
    //         if (mutation.type !== "childList") {
    //             continue;
    //         }
    //         // a line has been added or removed
    //         if (gutter.childElementCount < editarea.childElementCount) {
    //             let div = document.createElement("div");
    //             div.append(editarea.childElementCount);
    //             gutter.append(div);
    //         } else if (gutter.childElementCount > editarea.childElementCount) {
    //             gutter.removeChild(gutter.lastChild);
    //         }
    //     }
    // });
    // mutationObserver.observe(editarea, {
    //     childList: true
    // });
}

document.addEventListener("DOMContentLoaded", function () {
    const files = document.querySelectorAll('.fileListItem');
    files.forEach(function (element) {
        element.addEventListener('click', function (event) {
            files.forEach(function (element) {
                element.style.color = "";
                element.style.backgroundColor = "";
            });
            event.target.style.color = "white";
            event.target.style.backgroundColor = "#0749DA";

            fetch("/file/" + event.target.textContent)
                .then(response => response.text())
                .then(data => initEditor(data))
                .catch(error => console.error(error));
        });
    });
});

