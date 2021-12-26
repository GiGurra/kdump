package com.github.gigurra.kdump.internal.util

import scala.concurrent.duration.{Duration, FiniteDuration}
import scala.concurrent.{Await, ExecutionContext, Future}

object async {

  def run[T](expr: => T)(using ec: ExecutionContext): asyncOp[T] =
    asyncOp(Future(expr)(ec))

  def await[T](f: Future[T])(using ec: ExecutionContext): T =
    Await.result(f, Duration.Inf)

  def await[T](f: asyncOp[T])(using ec: ExecutionContext): T =
    await(f.op)

  case class asyncOp[T](op: Future[T]) {
    def join(using ec: ExecutionContext): T = this: T
  }

  given asyncOp2Val[T](using ec: ExecutionContext): Conversion[asyncOp[T], T] = x => await(x)
}
