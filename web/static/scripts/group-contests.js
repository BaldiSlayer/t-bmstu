function handleDivClick(url) {
    const target = event.target.tagName.toLowerCase();
    console.log(target);
    if (target !== 'a' && target !== 'i' && target !== 'img') {
        window.location.href = url;
    }
}