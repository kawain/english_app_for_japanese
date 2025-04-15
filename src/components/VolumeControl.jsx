function VolumeControl ({ volume, onVolumeChange }) {
  // スライダーが変更されたときに呼ばれる関数
  const handleChange = event => {
    // onVolumeChangeを通じてAppコンポーネントの状態を更新
    // event.target.value は文字列なので数値に変換
    onVolumeChange(parseInt(event.target.value, 10))
  }

  return (
    <div className='volume-control-wrapper'>
      <div className='volume-control'>
        <label htmlFor='volumeSlider'>音量: </label>
        <input
          type='range'
          id='volumeSlider'
          min='0'
          max='100'
          value={volume}
          onChange={handleChange}
        />
        <span id='volumeValue' className='volume-display'>
          {volume}%
        </span>
      </div>
    </div>
  )
}
export default VolumeControl
