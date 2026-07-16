export async function onRequestGet(ctx) {
    let path = ctx.request.url.toString();
    if (path.endsWith(".sig")) {
        return await ctx.env.ASSETS.fetch(ctx.request);
    }
    path = path + ".sig";
    const url = new URL(path);
    return await ctx.env.ASSETS.fetch(url);
}