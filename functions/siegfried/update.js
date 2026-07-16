export async function onRequestGet(ctx) {
    const url = new URL("/siegfried/update/update.json", ctx.request.url);
    return await ctx.env.ASSETS.fetch(url);
}