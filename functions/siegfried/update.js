export async function onRequestGet(ctx) {
    return await ctx.env.ASSETS.fetch("/siegfried/update/update.json");
}