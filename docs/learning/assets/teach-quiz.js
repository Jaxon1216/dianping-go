document.querySelectorAll("[data-quiz]").forEach((quiz) => {
  const buttons = Array.from(quiz.querySelectorAll("[data-answer]"));
  const feedback = quiz.querySelector("[data-feedback]");
  buttons.forEach((button) => {
    button.addEventListener("click", () => {
      buttons.forEach((choice) => {
        choice.setAttribute("aria-pressed", String(choice === button));
      });
      const correct = button.dataset.answer === quiz.dataset.correct;
      quiz.dataset.result = correct ? "correct" : "retry";
      feedback.textContent = correct
        ? quiz.dataset.success
        : quiz.dataset.retry;
    });
  });
});
