export const get = async (path: string) => {
  const res = await fetch(path, {
    method: "GET",
  });
  return await res.json();
};

export const post = async (path: string, body: any) => {
  const res = await fetch(path, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
  return await res.json();
};
