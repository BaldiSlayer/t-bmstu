function redirectToTask(taskId) {
    const currentURL = window.location.href;
    const contestIndex = currentURL.indexOf("contest/");
    const slashIndex = currentURL.indexOf("/", contestIndex + 8);
    const contestId = currentURL.substring(contestIndex + 8, slashIndex);
    const newURL = "/view/contest/" + contestId + "/problem/" + taskId;

    window.location.href = newURL;
}