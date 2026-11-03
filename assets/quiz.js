document.querySelectorAll("[data-quiz-question]").forEach((question) => {
  const feedback = question.querySelector("[data-feedback]");
  const options = [...question.querySelectorAll("[data-answer]")];

  options.forEach((option) => {
    option.addEventListener("click", () => {
      if (question.dataset.answered === "true") {
        return;
      }

      question.dataset.answered = "true";
      const isCorrect = option.dataset.answer === "correct";
      option.classList.add(isCorrect ? "correct" : "wrong");
      option.setAttribute("aria-pressed", "true");

      if (!isCorrect) {
        const correct = options.find(
          (candidate) => candidate.dataset.answer === "correct",
        );
        correct.classList.add("correct");
      }

      options.forEach((candidate) => {
        candidate.disabled = true;
      });

      feedback.textContent = isCorrect
        ? `正确。${question.dataset.explanation}`
        : `再想想。${question.dataset.explanation}`;
    });
  });
});
