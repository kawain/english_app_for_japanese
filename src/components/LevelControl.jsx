import { useAppContext } from '../App.jsx'

function LevelControl () {
  const {
    selectedLevel,
    handleLevelChange,
  } = useAppContext()

  // ラジオボタンの値が変更されたときに呼び出される関数
  const handleChange = event => {
    // 選択されたラジオボタンの値を渡す
    handleLevelChange(event.target.value)
  }

  return (
    <div className='level-control-wrapper'>
      <div className='level-select'>
        <label>
          <input
            type='radio'
            name='level'
            value='1'
            checked={selectedLevel === '1'}
            onChange={handleChange}
          />{' '}
          レベル1
        </label>
        <label>
          <input
            type='radio'
            name='level'
            value='2'
            checked={selectedLevel === '2'}
            onChange={handleChange}
          />{' '}
          レベル2
        </label>
        <label>
          <input
            type='radio'
            name='level'
            value='0'
            checked={selectedLevel === '0'}
            onChange={handleChange}
          />{' '}
          レベル1と2
        </label>
      </div>
    </div>
  )
}

export default LevelControl
