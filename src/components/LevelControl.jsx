function LevelControl ({ level, onLevelChange }) {
  // ラジオボタンの値が変更されたときに呼び出される関数
  const handleChange = event => {
    // 親コンポーネントに定義された onLevelChange を呼び出し、
    // 選択されたラジオボタンの値を渡す
    onLevelChange(event.target.value)
  }

  return (
    <div className='level-control-wrapper'>
      <div className='level-select'>
        <label>
          <input
            type='radio'
            name='level'
            value='1'
            // 受け取った level prop と比較して checked 状態を決定
            checked={level === '1'}
            // 変更時に handleChange 関数を呼び出す
            onChange={handleChange}
          />{' '}
          レベル1
        </label>
        <label>
          <input
            type='radio'
            name='level'
            value='2'
            checked={level === '2'}
            onChange={handleChange}
          />{' '}
          レベル2
        </label>
        <label>
          <input
            type='radio'
            name='level'
            value='9'
            checked={level === '9'}
            onChange={handleChange}
          />{' '}
          レベル1と2
        </label>
      </div>
    </div>
  )
}

export default LevelControl
