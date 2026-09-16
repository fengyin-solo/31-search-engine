async function getJSON(url) {
  const res = await fetch(url, { headers: { 'X-Auth-Token': 'dev-token' } });
  if (!res.ok) throw new Error(url + ' -> ' + res.status);
  return (await res.json()).data;
}

function tag(text) {
  const map = { active: '#16a34a', deleted: '#9ca3af', created: '#d97706', ready: '#16a34a', inactive: '#9ca3af' };
  const color = map[text] || '#6b7280';
  return `<span class="tag" style="background:${color}1a;color:${color}">${text || '-'}</span>`;
}

function renderStats(s) {
  const items = [
    ['文档', s.document_count], ['索引', s.index_count], ['分词器', s.analyzer_count],
    ['词条', s.term_count], ['倒排记录', s.posting_count], ['查询数', s.query_count],
  ];
  document.getElementById('stats').innerHTML = items.map(([label, num]) =>
    `<div class="card"><div class="num">${num}</div><div class="label">${label}</div></div>`).join('');
}

async function loadIndexes() {
  const data = await getJSON('/api/indexes');
  const sel = document.getElementById('index-select');
  sel.innerHTML = '';
  (data.items || []).forEach(idx => {
    const opt = document.createElement('option');
    opt.value = idx.id;
    opt.textContent = idx.name;
    sel.appendChild(opt);
  });
  const tbody = document.querySelector('#index-list');
  tbody.innerHTML = '';
  const stats = await getJSON('/api/stats/indexes');
  (stats.indexes || []).forEach(i => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${i.name}</td><td>${tag(i.status)}</td><td>${i.doc_count}</td><td>${i.term_count}</td>`;
    tbody.appendChild(tr);
  });
}

async function loadHot() {
  const data = await getJSON('/api/stats/queries');
  const tbody = document.querySelector('#hot-list');
  tbody.innerHTML = '';
  (data.top_hot_searches || []).forEach(h => {
    const tr = document.createElement('tr');
    tr.innerHTML = `<td>${h.term}</td><td>${h.count}</td>`;
    tbody.appendChild(tr);
  });
}

async function doSearch() {
  const indexID = document.getElementById('index-select').value;
  const q = document.getElementById('q').value.trim();
  if (!indexID) return;
  const data = await getJSON('/api/search?index_id=' + indexID + '&q=' + encodeURIComponent(q));
  const box = document.getElementById('results');
  if (!data.items || data.items.length === 0) {
    box.innerHTML = '<p style="color:#6b7280">未找到匹配结果</p>';
    return;
  }
  box.innerHTML = data.items.map(r => `
    <div style="border-bottom:1px solid #eef1f4;padding:10px 0">
      <div style="font-weight:600">${r.title} <span class="tag">${r.score.toFixed(2)}</span></div>
      <div style="color:#6b7280;font-size:13px">${r.body || ''}</div>
    </div>`).join('');
  await loadHot();
}

async function load() {
  const overview = await getJSON('/api/stats/overview');
  renderStats(overview);
  await loadIndexes();
  await loadHot();
}

load().catch(e => console.error(e));
