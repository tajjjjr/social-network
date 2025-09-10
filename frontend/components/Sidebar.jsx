if (!groupId || typeof groupId !== 'string' || groupId.length === 0) {
  setError('Invalid group ID');
  setLoading(false);
  return;
}

setLoading(true);
fetchData(groupId)
  .then(response => {
    // handle success
  })
  .catch(error => {
    setError('Failed to fetch data');
  })
  .finally(() => {
    setLoading(false);
  });