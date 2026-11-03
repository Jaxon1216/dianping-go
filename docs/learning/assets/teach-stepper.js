document.querySelectorAll("[data-stepper]").forEach((root) => {
  const source = root.querySelector("[data-stepper-frames]");
  if (!source) return;
  const frames = JSON.parse(source.textContent);
  let index = 0;
  const previous = root.querySelector('[data-move="previous"]');
  const next = root.querySelector('[data-move="next"]');
  const reset = root.querySelector('[data-move="reset"]');

  function write(selector, value) {
    const element = root.querySelector(selector);
    if (element) element.textContent = value;
  }

  function render() {
    const frame = frames[index];
    root.dataset.index = String(index);
    write("[data-step-number]", `${index + 1} / ${frames.length}`);
    write("[data-step-title]", frame.title);
    write("[data-active-code]", frame.code);
    write("[data-step-action]", frame.action);
    write("[data-step-input]", frame.input);
    write("[data-step-output]", frame.output);
    write("[data-step-reason]", frame.reason);

    const table = root.querySelector("[data-values]");
    table.replaceChildren();
    frame.values.forEach((row) => {
      const tr = document.createElement("tr");
      row.forEach((value) => {
        const td = document.createElement("td");
        td.textContent = value;
        tr.append(td);
      });
      table.append(tr);
    });
    if (!frame.values.length) {
      const tr = document.createElement("tr");
      const td = document.createElement("td");
      td.colSpan = 3;
      td.textContent = "还没有运行时变量或服务实例。";
      tr.append(td);
      table.append(tr);
    }

    write("[data-instance-label]", frame.count === null ? "服务实例尚未创建" : "S1 · 同一个 userService 实例");
    write("[data-instance-name]", frame.count === null ? "—" : "点评短信");
    write("[data-instance-count]", frame.count === null ? "—" : String(frame.count));
    write("[data-step-log]", frame.log || "尚无模拟发送输出。");
    previous.disabled = index === 0;
    next.disabled = index === frames.length - 1;
  }

  previous.addEventListener("click", () => {
    index = Math.max(0, index - 1);
    render();
  });
  next.addEventListener("click", () => {
    index = Math.min(frames.length - 1, index + 1);
    render();
  });
  reset.addEventListener("click", () => {
    index = 0;
    render();
  });
  root.dataset.mounted = "true";
  render();
});
